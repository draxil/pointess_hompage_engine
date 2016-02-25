package pointless_core;

import (
       "strings"
       "github.com/draxil/node_template"
       "log"
)      

const login_template = `
<form method='POST'>
      Password:<br>
      <input name=password type=password><br>
      Body:<br>
      <textarea rows=20 cols=80 name='page_data' id='page_data'></textarea>
      <br>
      <button>Save</button><br>
</form>`;

func ed_meat_el( current string ) ( *node_template.NodeTemplate ){
    meat_el, _ := node_template.Parse(strings.NewReader(login_template))
    page_data_el, _ := meat_el.FindFirst(`#page_data`);
    if( page_data_el != nil){
       page_data_el.ReplaceContentText(current)
    }
    log.Println( "Edit request for page  " + current )
    return meat_el
}