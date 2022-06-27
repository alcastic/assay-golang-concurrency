package main

import (
	"fmt"
	"time"
)

func main() {
	// 'The or-channel' pattern is uses to combining N 'done channels'
	// into just one 'done channel' which get close when any of the
	// original channels get closed
	var orChannel func(channels ...<-chan interface{}) <-chan interface{}
	orChannel = func(channels ...<-chan interface{}) <-chan interface{} {
		if len(channels) == 0 {
			return nil
		} else if len(channels) == 1 {
			return channels[0]
		}

		done := make(chan interface{})
		go func() {
			defer close(done)
			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-orChannel(append(channels[3:], done)...):
				}
			}
		}()
		return done
	}

	generateDoneChannel := func(d time.Duration) <-chan interface{} {
		done := make(chan interface{})
		go func() {
			time.Sleep(d)
			close(done)
		}()
		return done
	}

	done1 := generateDoneChannel(1 * time.Hour)
	done2 := generateDoneChannel(500 * time.Minute)
	done3 := generateDoneChannel(2 * time.Second)
	done4 := generateDoneChannel(3 * time.Minute)
	done5 := generateDoneChannel(5 * time.Minute)

	compositeDoneChannel := orChannel(done1, done2, done3, done4, done5)
	<-compositeDoneChannel
	fmt.Println("The compositeDoneChannel get close(child done channel have been closed)")
}
