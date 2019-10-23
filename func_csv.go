package starGo

import (
	"github.com/zhnxin/csvreader"
)

type Csv struct {
	reader *csvreader.Decoder
}

func NewCsvReader() *Csv {
	return &Csv{reader: csvreader.New()}
}

func (c *Csv) UnMarshalFile(file string, data interface{}) {
	err := c.reader.UnMarshalFile(file, data)
	if err != nil {
		ErrorLog("反序列化csv文件出错,错误信息:%v", err)
	}
}

func (c *Csv) UnMarshalFileWithHeader(file string, data interface{}, header []string) {
	err := c.reader.WithHeader(header).UnMarshalFile(file, &data)
	if err != nil {
		ErrorLog("反序列化csv文件出错,错误信息:%v", err)
	}
}
