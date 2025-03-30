package bench

import (
	"os"
	"testing"

	aliyun_oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/joho/godotenv"
	"github.com/tu6ge/oss-go"
)

func BenchmarkSelfUpload(b *testing.B) {
	client, _ := oss.NewWithEnv()
	obj := oss.NewObject("from_file.txt")

	f, _ := os.Open("./demofile.txt")

	defer f.Close()
	for b.Loop() {
		obj.File(f).ContentType("text/plain;charset=utf-8").Upload(&client)
	}
}

func BenchmarkAliyunUpload(b *testing.B) {
	godotenv.Load()

	// 读取环境变量
	key_id := os.Getenv("ALIYUN_KEY_ID")
	secret_id := os.Getenv("ALIYUN_KEY_SECRET")
	bucket_name := os.Getenv("ALIYUN_BUCKET")
	client, _ := aliyun_oss.New("https://oss-cn-shanghai.aliyuncs.com", key_id, secret_id)

	bucket, _ := client.Bucket(bucket_name)

	f, _ := os.Open("./demofile.txt")
	defer f.Close()
	for b.Loop() {
		bucket.PutObject("from_file.txt", f)
	}
}
