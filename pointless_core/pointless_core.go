package pointless_core;

import (
       "os"
       "io"
       "net/http"
       "log"
       "strings"
       "github.com/draxil/node_template"
       "code.google.com/p/go-sqlite/go1/sqlite3"
       "github.com/russross/blackfriday"
);


func Setup() (*phe_state){

     var dir string;
     if( len(os.Args) >= 2 ){
        dir = os.Args[1]
     } else {
        dir = "."
     }

     if( ! check_dir( dir ) ){
     	 return nil
     }  
  
     c, err := sqlite3.Open(dir + "/phe.db");

     if( err != nil ){
        log.Println("Problem with sqlite:" + err.Error() );
	return nil
     }    

     template, err := node_template.NodeTemplateFromFile( dir + "/skeleton.html")
     if( err != nil ){
     	 log.Println("Problem loading page skeleton: " + err.Error() )
	 return nil
     }

     state := phe_state{
     	   db: c,
	   skel: template,
	   dir: dir,	
	   first_run: false,
     };

     db_global_checks( &state )

     return &state
}

func check_dir( dir string ) ( bool ){
    log.Println( "examining " + dir )
    file, err := os.Stat(dir)
    if( err != nil ){
    	log.Println( "error with the selected main directory: " + err.Error() );
	return false;
    }
    if( ! file.IsDir() ){
    	log.Println( "error with the selected directory: not a directory!" );
	return false;
    }  
    return true;
}

type phe_state struct {
     db *sqlite3.Conn
     dir string
     skel *node_template.NodeTemplate
     first_run bool
};

func Page( w http.ResponseWriter, r * http.Request, state * phe_state ){

     edit_mode := r.FormValue("edit")
     check_save( r, state );

     if( state.first_run  ){
         first_run( w, r, state );
	 if( state.first_run ){
	     return;
	 }
     }

     meat, err := db_get_meat( state, r.URL.Path ) 

     var meat_el * node_template.NodeTemplate
     
     if( err == io.EOF ){
     	 meat_el = ed_meat_el("")
     } else if ( len(edit_mode) > 0 ){
       meat_el =  ed_meat_el( meat )
     } else {
       meat = string(blackfriday.MarkdownCommon( []byte(meat) ))
       meat_el, _ = node_template.Parse( strings.NewReader( meat ) );
    }
    render_page_with_meat( w, state, meat_el )
}

func render_page_with_meat( w http.ResponseWriter, state * phe_state, meat_el * node_template.NodeTemplate){

     page := state.skel.Copy(); 

     el, _ := page.Find("#meat")
     
     if( el != nil ){
     	 el.ReplaceContent( meat_el );
     }
     
     page.Render( w );
}

func check_save( r * http.Request, state * phe_state ){
     pd := r.FormValue("page_data")
     pwd := r.FormValue("password")
     log.Println( pwd );
     log.Println( crypt_site_password(pwd) );
     log.Println( db_site_password(state) );
     log.Println( crypt_site_password(pwd) == db_site_password(state) );
     if( len(pd) > 0  && crypt_site_password(pwd) == db_site_password(state) ){
         db_set_meat( state, r.URL.Path, pd );
     }    
}
