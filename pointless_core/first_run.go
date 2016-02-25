package pointless_core;

import( "net/http"
	"github.com/draxil/node_template"
	"strings" 
	"log"
);

func first_run( w http.ResponseWriter, r * http.Request, state * phe_state ){
     pwd := r.FormValue("password")
     if( len(pwd) > 0  ){
     	 db_set_config( state, "site_password", crypt_site_password(pwd) );
	 state.first_run = false;
	 return;
     }
     meat_el, err := node_template.Parse(strings.NewReader("<form method='POST'>First run, please set your password: <input type='password' name='password'><button>Save</button></form>"));

    if( err != nil ){
    	log.Println( err )
	return;
    }

    render_page_with_meat( w, state, meat_el )
}
