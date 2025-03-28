package oss

import (
	"testing"

	"github.com/tu6ge/oss-go/types"
)

func TestSecretEncryption(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<ListAllMyBucketsResult>
  <Owner>
    <ID>34773519</ID>
    <DisplayName>34773519</DisplayName>
  </Owner>
  <Buckets>
    <Bucket>
      <Comment></Comment>
      <CreationDate>2020-09-13T03:14:54.000Z</CreationDate>
      <ExtranetEndpoint>oss-cn-shanghai.aliyuncs.com</ExtranetEndpoint>
      <IntranetEndpoint>oss-cn-shanghai-internal.aliyuncs.com</IntranetEndpoint>
      <Location>oss-cn-shanghai</Location>
      <Name>aliyun-wb-kpbf3</Name>
      <Region>cn-shanghai</Region>
      <StorageClass>Standard</StorageClass>
    </Bucket>
    <Bucket>
      <Comment></Comment>
      <CreationDate>2016-11-05T13:10:10.000Z</CreationDate>
      <ExtranetEndpoint>oss-cn-shanghai.aliyuncs.com</ExtranetEndpoint>
      <IntranetEndpoint>oss-cn-shanghai-internal.aliyuncs.com</IntranetEndpoint>
      <Location>oss-cn-shanghai</Location>
      <Name>honglei123</Name>
      <Region>cn-shanghai</Region>
      <StorageClass>Standard</StorageClass>
    </Bucket>
  </Buckets>
</ListAllMyBucketsResult>`
	endpoint, _ := types.NewEndPoint("oss-qingdao")
	buckets := parser_xml(xml, endpoint)

	if buckets[0].name != "aliyun-wb-kpbf3" || buckets[1].name != "honglei123" {
		t.Error("parser xml failed")
	}
}
