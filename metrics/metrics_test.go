package metrics_test

import (
	"testing"

	"github.com/debeando/agent-mysql/metrics"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	m := metrics.Metric{}

	assert.True(t, m.Empty())

	m.AddTag(metrics.Tag{Name: "Test", Value: "test"})
	m.AddField(metrics.Field{Name: "Test", Value: 0})

	assert.False(t, m.Empty())
}

func TestCountTags(t *testing.T) {
	m := metrics.Metric{}

	assert.Equal(t, m.CountTags(), 0)

	m.AddTag(metrics.Tag{Name: "Test1", Value: "test1"})
	m.AddTag(metrics.Tag{Name: "Test2", Value: "test2"})

	assert.Equal(t, m.CountTags(), 2)
}

func TestCountFields(t *testing.T) {
	m := metrics.Metric{}

	assert.Equal(t, m.CountFields(), 0)

	m.AddField(metrics.Field{Name: "Test1", Value: "test1"})
	m.AddField(metrics.Field{Name: "Test2", Value: "test2"})

	assert.Equal(t, m.CountFields(), 2)
}
