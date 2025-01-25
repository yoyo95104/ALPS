package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/gin-gonic/gin"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"sync"
	"strconv"
)
type data struct{
	Name string
	Email string
	Grade int
	pass string
}
type response struct{
	Message string
	Status int
}
var (
	maxConcurrentRequests = 10 
	sem                     = make(chan struct{}, maxConcurrentRequests)
	wg                      sync.WaitGroup
)
func encrypt(text string, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))
	return ciphertext, nil
}

func generateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}
func handle(w http.ResponseWriter , r *http.Request){
	sem <- struct{}{}
	defer func() { <-sem }()
	time.Sleep(2 * time.Second)
	fmt.Fprintln(w , "request received")
	switch r.Method{
		case http.MethodPost:
			err := r.ParseForm()
			if err != nil {
				log.Fatal(err)
			}
			if r.FormValue("name") == " " && r.FormValue("grade") == " " {
				ch := make(chan bool)
				go verify(r.FormValue("email") , r.FormValue("pass") , ch)
				time.Sleep(1 * time.Second)
				if <-ch{
					fmt.Fprintln(w , true)
				}
		}else{
			ch := make(chan bool)
			grade, err := strconv.Atoi(r.FormValue("grade"))
			if err != nil {
				log.Println("Invalid grade value")
				return
			}
			go send(grade, r.FormValue("email"), r.FormValue("name"), r.FormValue("pass"), ch)
			grade , err = strconv.Atoi(r.FormValue("grade"))
			if err != nil {
				log.Println("Invalid grade value")
				return
			}
			go send(grade, r.FormValue("email"), r.FormValue("name"), r.FormValue("pass"), ch)
			time.Sleep(1 * time.Second)
			if <-ch{
				fmt.Fprintln(w , "Data inserted")
			}
		}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func main(){
	r := gin.Default()
	r.Run(":8080")
}

