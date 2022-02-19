package audio

import (
	"fmt"
	youtube "github.com/Genekoh/VikasBot/pkg/api"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"io"
	"log"
	"sync"
	"time"
)

type Song struct {
	title   string
	artist  string
	VideoId string
}

func NewSong(title string, artist string, videoId string) Song {
	return Song{
		title:   title,
		artist:  artist,
		VideoId: videoId,
	}
}

var VoiceInstances = make(map[string]*VoiceInstance)

type VoiceInstance struct {
	GuildId      string     // Guild Id
	InstanceDone chan error // Channel to signal voice instance is done
	//queueChannel chan Song
	Queue      []Song // Lists of Songs in Queue
	queueMutex sync.Mutex
	NowPlaying Song // Current Song Playing
	//streamDone      chan error                 // Channel signalling current song stream is done
	vc              *discordgo.VoiceConnection // Voice instance's voice connection
	currentStream   *dca.StreamingSession      //
	encodingSession *dca.EncodeSession
	play            chan bool
	pause           chan bool
	next            chan Song
	skip            chan bool
	isPlaying       bool
	isPaused        bool
	timeOut         *time.Timer
}

func GetVoiceInstance(guildId string, vc *discordgo.VoiceConnection) *VoiceInstance {
	vi := VoiceInstances[guildId]
	if vi == nil {
		vi = NewVoiceInstance(guildId, vc)
		VoiceInstances[guildId] = vi
	}
	fmt.Println("Got voice instance")

	return vi
}

func NewVoiceInstance(guildId string, vc *discordgo.VoiceConnection) *VoiceInstance {
	fmt.Println("creating voice instance")
	vi := VoiceInstance{
		GuildId:      guildId,
		InstanceDone: make(chan error),
		//queueChannel: make(chan Song),
		Queue: []Song{},
		vc:    vc,
		play:  make(chan bool),
		pause: make(chan bool),
		next:  make(chan Song),
		skip:  make(chan bool),
	}

	go establishInstanceLoop(&vi)
	fmt.Println("loop established")

	return &vi
}

func establishInstanceLoop(v *VoiceInstance) {
	fmt.Println("establishing loop")
	for {
		select {
		case nextSong := <-v.next:
			fmt.Println("received")
			v.NowPlaying = nextSong
			url, err := youtube.DownloadVideo(nextSong.VideoId)
			if err != nil {
				log.Println(err)
				continue
			}
			options := dca.StdEncodeOptions
			options.RawOutput = true
			options.Bitrate = 96
			options.Application = "lowdelay"

			encodingSession, err := dca.EncodeFile(url, options)
			if err != nil {
				log.Println(err)
				continue
			}
			v.encodingSession = encodingSession

			fmt.Println("sending stream")
			done := make(chan error)
			dca.NewStream(encodingSession, v.vc, done)
			go func() {
				v.play <- true
			}()
			go func() {
				err = <-done
				if err != nil && err != io.EOF {
					log.Println(err)
				}

				v.NextSong()
			}()

		case b := <-v.play:
			v.isPlaying = b
			if !b {
				v.timeOut = time.NewTimer(30 * time.Second)
			}
		}
	}
}

//func establishInstanceLoop(v *VoiceInstance) {
//	for {
//		select {
//		case b := <-v.play:
//			if b == v.isPlaying {
//				continue
//			}
//
//			v.isPlaying = b
//			if b == true {
//			} else {
//				v.timeOut = time.NewTimer(30 * time.Second)
//			}
//
//		case b := <-v.pause:
//			if b == v.isPaused {
//				continue
//			}
//
//			if b == true {
//				fmt.Println("pausing")
//				v.currentStream.SetPaused(true)
//				v.isPaused = true
//			} else {
//				fmt.Println("resuming")
//				v.currentStream.SetPaused(false)
//				v.isPlaying = false
//			}
//
//		case nextSong := <-v.queueChannel:
//			v.NowPlaying = nextSong
//			fmt.Println("downloading")
//			url, err := youtube.DownloadVideo(nextSong.VideoId)
//			if err != nil {
//				log.Println(err)
//				continue
//			}
//
//			fmt.Println("encoding options")
//			options := dca.StdEncodeOptions
//			options.RawOutput = true
//			options.Bitrate = 96
//			options.Application = "lowdelay"
//
//			fmt.Println("encoding")
//			encodingSession, err := dca.EncodeFile(url, options)
//			if err != nil {
//				log.Println(err)
//				continue
//			}
//
//			fmt.Println("sending stream")
//			done := make(chan error)
//			dca.NewStream(encodingSession, v.vc, done)
//			go func() {
//				err = <-done
//				if err != nil {
//					log.Println(err)
//				}
//
//				v.NextSong()
//			}()
//
//		case _ = <-v.timeOut.C:
//			v.DeleteInstance()
//			break
//
//		case err := <-v.InstanceDone:
//			v.DeleteInstance()
//			if err != nil {
//				log.Println(err)
//			}
//			break
//		}
//	}
//}

func (v *VoiceInstance) NextSong() {
	fmt.Println("next song called")
	if v.encodingSession != nil {
		fmt.Println("cleaning up")
		v.encodingSession.Cleanup()
	}

	if len(v.Queue) == 0 {
		go func() {
			v.play <- false
		}()
		return
	}

	nextSong := v.Queue[0]
	v.QueueRemove(0)
	go func() {
		v.next <- nextSong
	}()
}

func (v *VoiceInstance) QueueAdd(song Song) {
	fmt.Println("before go routine")
	go func() {
		v.queueMutex.Lock()
		v.Queue = append(v.Queue, song)
		v.queueMutex.Unlock()

		if v.isPlaying == false {
			fmt.Println("Next Song called")
			v.NextSong()
		}
	}()
}

func (v *VoiceInstance) QueueRemove(i int) {
	v.queueMutex.Lock()
	defer v.queueMutex.Unlock()

	newQueue := []Song{}
	if len(v.Queue) != 0 && i <= len(v.Queue) {
		newQueue = append(v.Queue[:i], v.Queue[i+1:]...)
	} else {
		newQueue = []Song{}
	}

	v.Queue = newQueue
}

func (v *VoiceInstance) QueueClear() {
	v.queueMutex.Lock()
	defer v.queueMutex.Unlock()

	v.Queue = []Song{}
}

func (v *VoiceInstance) DeleteInstance() {
	err := v.vc.Disconnect()
	if err != nil {
		log.Println(err)
		return
	}
	delete(VoiceInstances, v.GuildId)
	*v = VoiceInstance{}
}
