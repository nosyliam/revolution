package detect

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/nosyliam/revolution/pkg/config"
	"hash/crc32"
	"image"
	"sync"
)

type BuffType int

const (
	Haste BuffType = iota
	HastePlus
	HasteCoconut
	BearMorph
	Focus
)

type Buff struct {
	Type  BuffType
	Value int
}

type BuffList []Buff

func (b *BuffList) Hash() string {
	h := crc32.New(crc32.MakeTable(crc32.IEEE))

	for _, buff := range *b {
		var buf bytes.Buffer
		_ = binary.Write(&buf, binary.LittleEndian, buff.Type)
		_ = binary.Write(&buf, binary.LittleEndian, int32(buff.Value))
		h.Write(buf.Bytes())
	}

	return fmt.Sprintf("%08x", h.Sum32())
}

type BuffDetector struct {
	mu              sync.Mutex
	hourlyHistogram map[BuffType][]int
	interrupts      map[<-chan struct{}]bool
	buffs           []Buff
	tick            int
	settings        *config.Object[config.Settings]
}

func NewBuffDetector(settings *config.Object[config.Settings]) *BuffDetector {
	return &BuffDetector{
		interrupts:      make(map[<-chan struct{}]bool),
		hourlyHistogram: make(map[BuffType][]int),
	}
}

func (b *BuffDetector) MoveSpeed() float64 {
	speed := *config.Concrete[float64](b.settings, "player.moveSpeed")
	for _, buff := range b.buffs {
		switch buff.Type {
		case BearMorph:
			speed += 8 // Bear morph must be added first
		}
	}
	for _, buff := range b.buffs {
		switch buff.Type {
		case Haste:
		case HastePlus:
		case HasteCoconut:

		}
	}
	return speed

}

// Combine generates a channel which combines the global macro interrupt with an interrupt for speed buff changes
func (b *BuffDetector) Combine(interrupt <-chan struct{}) <-chan struct{} {
	b.mu.Lock()
	defer b.mu.Unlock()
	waiter := make(chan struct{})
	b.interrupts[waiter] = true
	combined := make(chan struct{})
	go func() {
		select {
		case <-interrupt:
			b.mu.Lock()
			delete(b.interrupts, waiter)
			b.mu.Unlock()
			combined <- struct{}{}
		case <-waiter:
			b.mu.Lock()
			delete(b.interrupts, waiter)
			b.mu.Unlock()
			combined <- struct{}{}
		}
	}()
	return combined
}

func (b *BuffDetector) Tick(screenshot *image.RGBA) {

}
