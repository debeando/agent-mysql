package metrics

// Tag for metric.
type Tag struct {
	Name  string
	Value string
}

// Field is a collection for specific metric value.
type Field struct {
	Name  string
	Value interface{}
}

// Metric is a collection with many Tags and Fields.
type Metric struct {
	Measurement string
	Tags        []Tag
	Fields      []Field
}

// Metrics of a colletion of metric.
type Metrics []Metric

// Verify Metric is not empty on Tags or Fields.
func (m *Metric) Empty() bool {
	return ((*m).CountTags() == 0) || ((*m).CountFields() == 0)
}

// Add is aggregator for metric into metrics.
func (m *Metric) AddTag(i Tag) {
	(*m).Tags = append((*m).Tags, i)
}

// Add is aggregator for field into metrics.
func (m *Metric) AddField(i Field) {
	(*m).Fields = append((*m).Fields, i)
}

// Count Tags on metric.
func (m *Metric) CountTags() int {
	return len((*m).Tags)
}

// Count Fields on metric.
func (m *Metric) CountFields() int {
	return len((*m).Fields)
}

// Convert Tags to Map.
func (m *Metric) TagsToMap() map[string]string {
	tags := make(map[string]string)
	for _, tag := range (*m).Tags {
		tags[tag.Name] = tag.Value
	}

	return tags
}

// Convert Fields to Map.
func (m *Metric) FieldsToMap() map[string]interface{} {
	fields := make(map[string]interface{})
	for _, field := range (*m).Fields {
		fields[field.Name] = field.Value
	}

	return fields
}

// Count all metrics in metrics.
func (m *Metrics) Count() int {
	return len(*m)
}

// Reset the metric metrics.
func (m *Metrics) Reset() {
	*m = (*m)[:0]
}

// Add is aggregator for metric into metrics.
func (m *Metrics) Add(i Metric) {
	if !i.Empty() {
		*m = append(*m, i)
	}
}
