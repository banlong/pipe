package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"os/exec"
)

func main() {
	http.HandleFunc("/hijack", hijackHandler)
	http.HandleFunc("/pipe", testPipe)
	http.HandleFunc("/viewyt", Tran)
	http.HandleFunc("/stream", Stream)
	log.Println("HTTP Server...9092")
	log.Fatal(http.ListenAndServe(":9092", nil))

}

func testPipe(w http.ResponseWriter, r *http.Request)  {
	//Use this to communicate with client through pipe:
	// 1. Receive data from client
	// 2. Write message to client
	// 3. pipe is not close
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	log.Printf("Received: %s\n", b)
	w.Write([]byte("hello\n"))
	w.Write([]byte("this is me\n"))
}

func hijackHandler(w http.ResponseWriter, r *http.Request)  {
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	//Hijack http connection into TCP connection
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Don't forget to close the connection:
	defer conn.Close()

	//Write to the client
	bufrw.WriteString("Now we're speaking raw TCP. Say hi: \n")
	bufrw.Flush()

	//Read from the client
	s, err := bufrw.ReadString('\n')
	if err != nil {
		log.Printf("error reading string: %v", err)
		return
	}
	fmt.Fprintf(bufrw, "You said: %q\nBye.\n", s)
	bufrw.Flush()
}

func encodeAudioHandler(w http.ResponseWriter, req *http.Request) {

	streamLink := exec.Command("youtube-dl", "-f", "140", "-g","--no-check-certificate","https://www.youtube.com/watch?v=MpQR4IsdKSs")
	out, err := streamLink.Output()

	//display out put
	log.Println("OUT:", string(out))

	if err != nil {
		log.Println("Link err >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		log.Println(err.Error())
	}

	resp, err := http.Get(string(out))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	//set the output of ffmpeg to stdout
	cmdFF := exec.Command("ffmpeg", "-i", "pipe:0", "-acodec", "libmp3lame", "-f", "mp3", "-")
	//pipe the response content to stdin (input for ffmpeg)
	cmdFF.Stdin = resp.Body

	/////////result of cmdFF.RUN send to  file
	// open the out file for writing
//	outfile, err := os.Create("out.mp4")
//	if err != nil {
//		panic(err)
//	}
//	defer outfile.Close()
//	cmdFF.Stdout = outfile =
//	////////////////////////////////////////////////////////
//
//	//link output of the cmdFF to response writer
	cmdFF.Stdout = w
	if err := cmdFF.Run(); err != nil {
		log.Println("runFF err >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		log.Println(err.Error())
	}
}

func Tran(w http.ResponseWriter, req *http.Request) {
	cmd := fmt.Sprintln("ffmpeg" + " -i " + " videos/sample.mp4" +
	" -y -vcodec libx264 -preset ultrafast -threads 0 -c:a aac -strict -2 " + " out.mp4")
	cmdFF := exec.Command("bash", "-c", cmd)
	err := cmdFF.Run()
	if(err != nil){
		log.Println(err.Error())
	}
}

func Stream(w http.ResponseWriter, req *http.Request){
	cmd := fmt.Sprintln("ffmpeg -i videos/sample.mp4 -f mp4 pipe:2")
	cmdFF := exec.Command("bash", "-c", cmd)
	cmdFF.Stdout = w
	err := cmdFF.Run()
	if(err != nil){
		log.Println(err.Error())
	}

	//vReader := bytes.NewReader(cmdFF.Stdout)
	//http.ServeContent(w, req, "sample.mp4", time.Now(), cmdFF.Stdout)
}