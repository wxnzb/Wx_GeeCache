package geecache
import(
   "testing"
   "reflect"
)
func Test_Getter(t *testing.T){
     var f=GetterFunc(key string)([]byte,error){
		return []byte(key),nil
	 }
	 expect:=[]byte("key")
	 if v,_:=f.Get("key");!reflect.DeepEqual(v,expect){
		 t.Errorf("got %v, expect %v",v,expect)
	 }
}