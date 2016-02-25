package pointless_core;

import ( 
       "io" 
       "log"
       "code.google.com/p/go-sqlite/go1/sqlite3"
)


/// TODO: needs mutex on db


func db_get_meat (s * phe_state, path string ) (string, error){
	 db := s.db

	 stmt, err := db.Query( "select meat from page where path=?", path );
	 
	 if( err != nil ){

	     if( err == io.EOF ){
	     	 return "",  io.EOF;
	     }

	     log.Println("error with page: " + err.Error() );
	     return "Error.", nil;
         }	

	 var meat string;

	 stmt.Scan(&meat)

	 return meat, nil; 
}

func check_page_table( db * sqlite3.Conn ){

	if( ! table_exists( db, "page") ){
	    db.Exec("create table page ( path varchar(200) primary key, meat TEXT )");
	}
}
func db_set_meat( s * phe_state, path string, meat string ) {
	 db := s.db
	 check_page_table( db );
	 db.Exec("replace into page values(?, ?)", path, meat )
}

func table_exists( db * sqlite3.Conn, table string ) (bool){
	_, err := db.Query( "SELECT name FROM sqlite_master WHERE type='table' AND name=?", table)
	if( err != nil && err != io.EOF ){
	    log.Println(err);
	}
	if( err != nil && err == io.EOF ){
	    return false
	}
	return true
}

func db_global_checks( s * phe_state ){
     db := s.db
     if( ! table_exists( db, "cfg" ) ){
     	   log.Println("no config table...");
	   err := db.Exec("create table cfg ( setting varchar(200) primary key, value varchar(200) )");
	   if( err != nil ){
	       log.Println(err);
	   }
     } else{
	   log.Println("have config table...");
     }
     check_page_table( db );
     if( db_site_password( s ) == "" ){
     	 s.first_run = true;
     }
}
func db_site_password( s * phe_state ) (string){
     return db_config_setting( s, "site_password");
}

func db_config_setting( s * phe_state, setting string ) (string){
     	 db := s.db
	 stmt, err := db.Query( "select value from cfg where setting=?", setting )
	 if( err == nil ){
	     	 var val string
	 	 stmt.Scan(&val)
		 return val
         }	
	 if( err != io.EOF ){
	     log.Println(err);
	 }
	 return ""
}

func db_set_config( s * phe_state, setting string, value string ){
     db := s.db;
     db.Exec("replace into cfg values( ?, ?)", setting, value )
}