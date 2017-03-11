package ymqd

import "testing"

func Test_generateGuid(t *testing.T) {

	topic, _ := NewTopic(nil)
	for i := 0; i < 1000; i++ {
		topic.GenerateID()
	}
}

func Benchmark_generateGuid(t *testing.B) {
	topic, _ := NewTopic(nil)
	for i := 0; i < 10000; i++ {
		topic.GenerateID()
	}
}
