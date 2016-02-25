package main;

import (
       "net/http"
       "github.com/draxil/pointless_homepage_engine/pointless_core"
);


func main(){
     state := pointless_core.Setup();
     if( state == nil ){
     	 return
     }
     http.HandleFunc( "/",
		   func( w http.ResponseWriter, r *http.Request){
		       pointless_core.Page( w, r, state );
		   });
     http.ListenAndServe(":4000", nil)
}

