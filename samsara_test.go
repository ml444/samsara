package samsara

import (
	"bytes"
	"github.com/ml444/samsara/core"
	"github.com/ml444/samsara/entity"
	"github.com/ml444/samsara/publish"
	"github.com/ml444/samsara/subscribe"
	"github.com/ml444/samsara/utils"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewSamsara(t *testing.T) {
	ringBuffer := core.NewRingBuffer(1024, &entity.FactoryEntity{})

	type args struct {
		ringBufferSize int64
		eventFactory   core.IEntityFactory
	}
	tests := []struct {
		name string
		args args
		want *Samsara
	}{
		{
			name: "",
			args: args{
				ringBufferSize: 1024,
				eventFactory:   &entity.FactoryEntity{},
			},
			want: &Samsara{
				ringBuffer: ringBuffer,
				scheduler:  core.NewScheduler(ringBuffer),
				isDone:     &utils.AtomicBool{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSamsara(tt.args.ringBufferSize, tt.args.eventFactory); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSamsara() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSamsara_Get(t *testing.T) {
	type args struct {
		sequence int64
	}
	tests := []struct {
		name string
		args args
		want core.IEntity
	}{
		{
			name: "",
			args: args{1},
			want: &entity.Entity{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSamsara(1024, &entity.FactoryEntity{})
			if got := s.Get(tt.args.sequence); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSamsara_GetBufferSize(t *testing.T) {
	type fields struct {
		ringBufferSize int64
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name:   "",
			fields: fields{1024},
			want:   1024,
		},
		{
			name:   "",
			fields: fields{10240},
			want:   10240,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSamsara(tt.fields.ringBufferSize, &entity.FactoryEntity{})
			if got := s.GetBufferSize(); got != tt.want {
				t.Errorf("GetBufferSize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func NewTestSamsara() *Samsara {
	s := NewSamsara(102400, &entity.FactoryEntity{})
	producer := publish.NewSinglePublisher(s.GetScheduler(), publish.NewSinglePublishStrategy(10*time.Microsecond))
	s.SetPublisher(producer)
	consumer := subscribe.NewSimpleSubscriber(s.GetScheduler(), subscribe.NewSingleSubscribeStrategy(10*time.Microsecond))
	s.SetSubscriber(consumer)
	s.Start()
	return s
}
func NewTestSamsaraWithMultiPublisher() *Samsara {
	s := NewSamsara(102400, &entity.FactoryEntity{})
	producer := publish.NewMultiPublisher(s.GetScheduler(), publish.NewSinglePublishStrategy(10*time.Microsecond))
	s.SetPublisher(producer)
	consumer := subscribe.NewSimpleSubscriber(s.GetScheduler(), subscribe.NewSingleSubscribeStrategy(10*time.Microsecond))
	println(consumer)
	s.SetSubscriber(consumer)
	s.Start()
	return s
}

func TestSamsara_GetCursor(t *testing.T) {
	tests := []struct {
		name   string
		pubNum int
		want   int64
	}{
		{
			name:   "001",
			pubNum: 0,
			want:   -1,
		},
		{
			name:   "002",
			pubNum: 10,
			want:   9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewTestSamsara()
			for i := 0; i < tt.pubNum; i++ {
				_ = s.Publish(&entity.Entity{Data: []byte{}})
			}
			if got := s.GetCursor(); got != tt.want {
				t.Errorf("GetCursor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSamsara_GetRingBuffer(t *testing.T) {
	type fields struct {
		ringBufferSize int64
	}
	tests := []struct {
		name   string
		fields fields
		want   *core.RingBuffer
	}{
		{
			name:   "test",
			fields: fields{ringBufferSize: 1024},
			want:   core.NewRingBuffer(1024, &entity.FactoryEntity{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSamsara(tt.fields.ringBufferSize, &entity.FactoryEntity{})
			if got := s.GetRingBuffer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRingBuffer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSamsara_HasBlocking(t *testing.T) {
	type fields struct {
		ringBufferSize int64
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "test",
			fields: fields{ringBufferSize: 1024},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSamsara(tt.fields.ringBufferSize, &entity.FactoryEntity{})
			if got := s.HasBlocking(); got != tt.want {
				t.Errorf("HasBlocking() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSamsara_IsDone(t *testing.T) {
	tests := []struct {
		name       string
		isShutdown bool
		want       bool
	}{
		{name: "test_running", isShutdown: false, want: false},
		{name: "test_done", isShutdown: true, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewTestSamsara()
			if tt.isShutdown {
				s.Shutdown()
			}
			if got := s.IsDone(); got != tt.want {
				t.Errorf("IsDone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSamsara_Publish(t *testing.T) {
	tests := []struct {
		name   string
		pubNum int
		want   int64
	}{
		{name: "001", pubNum: 0, want: -1},
		{name: "002", pubNum: 10, want: 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewTestSamsara()
			for i := 0; i < tt.pubNum; i++ {
				_ = s.Publish(&entity.Entity{Data: []byte{}})
			}
			if got := s.GetCursor(); got != tt.want {
				t.Errorf("GetCursor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSamsara_SetConsumer(t *testing.T) {
	tests := []struct {
		name   string
		pubNum int
		want   int64
	}{
		{name: "001", pubNum: 0, want: -1},
		{name: "002", pubNum: 10, want: 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSamsara(1024, &entity.FactoryEntity{})
			producer := publish.NewSinglePublisher(s.GetScheduler(), publish.NewSinglePublishStrategy(10*time.Microsecond))
			s.SetPublisher(producer)
			consumer := subscribe.NewSimpleSubscriber(s.GetScheduler(), subscribe.NewSingleSubscribeStrategy(10*time.Microsecond))
			s.SetSubscriber(consumer)
			s.Start()
			for i := 0; i < tt.pubNum; i++ {
				_ = s.Publish(&entity.Entity{Data: []byte{}})
			}
			time.Sleep(100 * time.Millisecond)
			if got := consumer.GetSequence().Get(); got != tt.want {
				t.Errorf("GetCursor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSamsara_SetPublisher(t *testing.T) {
	s := NewSamsara(1024, &entity.FactoryEntity{})
	publisher := publish.NewSinglePublisher(s.GetScheduler(), publish.NewSinglePublishStrategy(10*time.Microsecond))
	s.SetPublisher(publisher)
	if !reflect.DeepEqual(s.publisher, publisher) {
		t.Errorf("SetPublisher is error")
	}
}

func TestSamsara_Shutdown(t *testing.T) {
	s := NewTestSamsara()
	s.Shutdown()
	if !s.IsDone() {
		t.Errorf("Shutdown is failed.")
	}
}

func TestSamsara_Start(t *testing.T) {

}

func TestSamsara_StopSubscribers(t *testing.T) {
	s := NewTestSamsara()
	for i := 0; i < 100; i++ {
		_ = s.Publish(&entity.Entity{Data: []byte{}})
	}
	time.Sleep(1000 * time.Millisecond)
	s.StopSubscribers()
	for i := 0; i < 11; i++ {
		_ = s.Publish(&entity.Entity{Data: []byte{}})
	}
	if s.GetCursor() == 110 {
		for _, sub := range s.subscriberList {
			if sub.GetSequence().Get() != 99 {
				t.Errorf("StopSubscribers is failed.")
			}
		}
	}
}

func TestSamsara_StopPublisher(t *testing.T) {
	s := NewTestSamsara()
	for i := 0; i < 100; i++ {
		err := s.Publish(&entity.Entity{Data: []byte{}})
		if err != nil {
			t.Errorf("Err: %v\n", err)
		}
	}
	s.StopPublisher()
	for i := 0; i < 11; i++ {
		err := s.Publish(&entity.Entity{Data: []byte{}})
		if err == nil {
			t.Errorf("err must be not nil.")
		}
	}
}

func TestLog(t *testing.T) {

	//f, err := os.Create("binlog_pprof1.pprof")
	//if err != nil {
	//	log.Fatal("could not create CPU profile: ", err)
	//}
	//if err := pprof.StartCPUProfile(f); err != nil {
	//	log.Fatal("could not start CPU profile: ", err)
	//}
	//defer pprof.StopCPUProfile()

	s := NewTestSamsaraWithMultiPublisher()
	buf := bytes.NewBufferString(`INF exec.go:github.com/ml444/samsara/impl.RawQuery:454 [DBPROXY_EXEC_READ] [4]  SELECT * FROM id_config WHERE (tb_name = 'm_member') AND (deleted_at=0 OR deleted_at IS NULL) LIMIT 1 [1 rows]`)
	buf.WriteString("\n")
	wg := sync.WaitGroup{}
	startTime := time.Now()
	for j := 0; j < 1000; j++ {
		wg.Add(1)
		go func(n int) {
			defer func() {
				wg.Done()
				//fmt.Println("===>", n)
			}()
			for i := 0; i < 200; i++ {
				err := s.Publish(&entity.Entity{Data: buf.Bytes()})
				if err != nil {
					t.Errorf("Err: %v\n", err)
				}
			}
		}(j)
	}
	wg.Wait()

	t.Log("===> speed time:", time.Now().Sub(startTime).Milliseconds())
	time.Sleep(1 * time.Second)
	subscribe.FileFlush()
	time.Sleep(1 * time.Second)
	s.Shutdown()
}

func TestMod(t *testing.T) {
	var a int64 = 10241024123
	var mask int64 = 1023
	t.Log("===>", a&mask)
}

func TestFor(t *testing.T) {
	var a = 1
	get := func() int {
		a++
		return a
	}
	for i := get(); i < 10; {
		t.Log(i)
		time.Sleep(time.Second)
	}
}