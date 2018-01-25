package main
import (
"net/http"
	"fmt"
	"sync"
	"time"
)
type mengsSever struct{}
type Args struct {
	A string
	UID string
	Stime time.Time
	Room string
}
type RoomList struct {
	roomList map[string]string
	Lock sync.Mutex
}

type History struct {
	history map[string]string
	Lock sync.Mutex
}

type NameID struct {
	nameid map[string]string
	Lock sync.Mutex
}

type RoomCount struct {
	roomcount map[string]int
	Lock sync.Mutex
}

type HistoryTime struct {
	historytime map[string]time.Time
	Lock sync.Mutex
}
var roomList RoomList
var history History
var nameID NameID
var roomCount RoomCount
var historyTime HistoryTime
var message string
var msg []byte
var flag=false
var arg1 Args

func (rl RoomList)Get(k string) string{
	rl.Lock.Lock()
	defer rl.Lock.Unlock()
	return rl.roomList[k]
}
func (rl RoomList)Set(k,v string) {
	rl.Lock.Lock()
	defer rl.Lock.Unlock()
	rl.roomList[k]=v
}
/*func (rl RoomList)Delete(k string)  {
	rl.Lock.Lock()
	defer rl.Lock.Unlock()
	delete(rl.roomList,k)
}*/

func (hs History)Get(k string) string{
	hs.Lock.Lock()
	defer hs.Lock.Unlock()
	return hs.history[k]
}
func (hs History)Set(k,v string) {
	hs.Lock.Lock()
	defer hs.Lock.Unlock()
	hs.history[k]=v
}

/*func (nid NameID)Get(k string) string{
	nid.Lock.Lock()
	defer nid.Lock.Unlock()
	return nid.nameid[k]
}
func (nid NameID)Set(k,v string) {
	nid.Lock.Lock()
	defer nid.Lock.Unlock()
	nid.nameid[k]=v
}

func (rc RoomCount)Get(k string) int{
	rc.Lock.Lock()
	defer rc.Lock.Unlock()
	return rc.roomcount[k]
}
func (rc RoomCount)Set(k string,v int) {
	rc.Lock.Lock()
	defer rc.Lock.Unlock()
	rc.roomcount[k]=v
}*/

func (ht HistoryTime)Get(k string) time.Time{
	ht.Lock.Lock()
	defer ht.Lock.Unlock()
	return ht.historytime[k]
}
func (ht HistoryTime)Set(k string,v time.Time) {
	ht.Lock.Lock()
	defer ht.Lock.Unlock()
	ht.historytime[k]=v
}

/*func changeToInt(s string) int  {
	int,_:=strconv.Atoi(s)
	return int
}*/
func CheckRoom()  {
	var subTime time.Duration
	for true {
		for roomID := range historyTime.historytime{
			subTime =time.Now().Sub(historyTime.Get(roomID))
			if subTime>30*time.Second {
				fmt.Println(subTime.Seconds(),roomID)
				flag=true
				message="#"+roomID+"&"+"龖"+"|"
				/*RoomCount[roomID]=-1*/
				fmt.Println(message)
				time.Sleep(time.Second)
				delete(history.history,roomID)
				delete(historyTime.historytime,roomID)
			}
		}
		time.Sleep(1*time.Second)
	}
}
func (h mengsSever) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	lenth:= req.ContentLength
	var bodySlc []byte = make([]byte, lenth)
	s:=req.Method
	/*fmt.Println(s)*/
	if s=="POST" {
		bodyLen,readErr := req.Body.Read(bodySlc)
		arg1.Room=req.FormValue("Room")
		arg1.A =  req.FormValue("Message")
		arg1.UID = req.FormValue("Uid")
		if readErr != nil {
			fmt.Println("read body error")
		}else{
			fmt.Println("the body has ",bodyLen," bytes")
		}
		/*message:= (string(bodySlc))*/
		if arg1.UID!= "###历史###"{
			flag=true
			t:=time.Now()
			record:= "\n"+history.Get(arg1.Room)+arg1.UID + " says :"+ arg1.A+"\n"+t.Format(time.RFC850)+"\n"
			history.Set(arg1.Room,record)
			fmt.Println(history.Get(arg1.Room))
			historyTime.Set(arg1.Room,t)
			hs:=history.Get(arg1.Room)
			message = hs+"#"+arg1.Room+"&"+arg1.UID +"|"+" says :"+ arg1.A+"\n"+t.Format(time.RFC850)
			/*fmt.Println("the body is:", message)*/
			fmt.Println(message)
		}else {
			t:=time.Now()
			record:= "\n"+history.Get(arg1.Room)+arg1.A+"    "+t.Format(time.RFC850)+"\n"
			history.Set(arg1.Room,record)
			for room:=range history.history{
				if room==arg1.Room {
					History:=history.Get(room)
					w.Write([]byte(History))
				}
			}
		}

	}else {
		if flag==true {
			w.Write([]byte(message))
			time.Sleep(11*time.Millisecond)
			flag=false
		}else {
			var rooms string
			for room:=range history.history {
				rooms+="["+room+"]"
			}
			h:="齾"+"#"+rooms
			w.Write([]byte(h))
		}
	}



}

func main() {
	roomList.roomList=make(map[string]string)
	history.history=make(map[string]string)
	nameID.nameid=make(map[string]string)
	roomCount.roomcount=make(map[string]int)
	historyTime.historytime=make(map[string]time.Time)
	go CheckRoom()
	var s mengsSever
	http.Handle("/",s )
	http.ListenAndServe(":12121", nil)
	}
/*



type helloHandler struct{}

func (h *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func main() {
	http.Handle("/", &helloHandler{})
	http.ListenAndServe(":12345", nil)
}*/