package model

type Visitor interface {
	BeforeStruct(data *MetaData)
	AfterStruct(data *MetaData)

	BeforeArray(data *MetaData)
	AfterArray(data *MetaData)

	BeforeSlice(data *MetaData)
	AfterSlice(data *MetaData)

	NamedField(data *MetaData)
	IndexedField(data *MetaData)
}
