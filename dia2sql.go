package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "compress/gzip"
    "strings"
    "github.com/antchfx/xmlquery"
)
var table_name_str, table_field_name_str,table_field_datatype_str,primary_keys_list_str,unique_field_names_list_str string
var is_verbose_output=false  
func main() {
    if (len(os.Args[1:])==0){
	    fmt.Println("Usage..: dia2sql [-v] your_dia_file.dia ")
	    return
    }
    argsWithoutProgramName := os.Args[1:]
    if (argsWithoutProgramName[0]=="-v") {
	    is_verbose_output=true
    }
    filename:=argsWithoutProgramName[len(argsWithoutProgramName)-1]
    if (len(filename)==0) {
	  fmt.Println("Usage..: dia2sql [-v] your_dia_file.dia")
	  return  
    }
    file, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    gz, err := gzip.NewReader(file)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	defer gz.Close()

	scanner := bufio.NewScanner(gz)
    // optionally, resize scanner's capacity for lines over 64K, see next example
	xml_contents:=""
    for scanner.Scan() {
        xml_contents=xml_contents+strings.Trim(scanner.Text()," ")
	
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
    doc, err := xmlquery.Parse(strings.NewReader(xml_contents))
    database_tables_as_xml_nodes_list,_ := xmlquery.QueryAll(doc, "//dia:diagram/dia:layer/dia:object[@type=\"Database - Table\"]")
    
    for index:=0; index<len(database_tables_as_xml_nodes_list); index++ {
	    
	    
	         database_table_xml_node:=database_tables_as_xml_nodes_list[index]
	         path:="//dia:attribute[@name=\"name\"]/dia:string"
	         database_table_name_as_xml_node,_:=xmlquery.Query(database_table_xml_node, path)
		 //fmt.Println(database_table_name)
		 table_name_str=strings.Trim(database_table_name_as_xml_node.InnerText(),"#")
		 fmt.Println("\n-- %%Queries for table ",table_field_name_str)
		if (is_verbose_output) { fmt.Println("table name ",index+1, " ",table_name_str ) }
		 table_sql_query:="create table if not exists "+table_name_str+" ("
		 
		 path="//dia:attribute[@name=\"attributes\"]/dia:composite[@type=\"table_attribute\"]"
		 database_table_attributes_as_xml_nodes_list,_:=xmlquery.QueryAll(database_table_xml_node, path)
		 
		 for index1:=0; index1<len(database_table_attributes_as_xml_nodes_list); index1++ {
			if (is_verbose_output) { fmt.Println("=========") }
			 table_field_as_xml_node:=database_table_attributes_as_xml_nodes_list[index1]
			 path_attribs:="//dia:attribute[@name=\"name\"]/dia:string"
			 table_field_name,_:=xmlquery.Query(table_field_as_xml_node,path_attribs)
			 table_field_name_str=strings.Trim(table_field_name.InnerText(),"#")
			if (is_verbose_output) { fmt.Println(table_field_name_str)}
			 path_attribs="//dia:attribute[@name=\"type\"]/dia:string"
			 table_field_datatype,_:=xmlquery.Query(table_field_as_xml_node,path_attribs)
			 table_field_datatype_str=strings.Trim(table_field_datatype.InnerText(),"#")
			 if (table_field_datatype_str=="") {
				 fmt.Println("error in table "+table_name_str+". the field "+table_field_name_str+" has no datatype defined.\nEdit the diagram file and fill the datatype for the field "+table_field_name_str)  
				 return
			 }
			if (is_verbose_output) { fmt.Println(table_field_datatype_str)}
			 table_sql_query+=table_field_name_str+" " +table_field_datatype_str+"," 
			 
			 
			 path_attribs="//dia:attribute[@name=\"primary_key\"]/dia:boolean[@val=\"true\"]"
			 table_field_is_primary_key,_:=xmlquery.Query(table_field_as_xml_node,path_attribs)
			 if(table_field_is_primary_key!=nil) {
				 primary_keys_list_str+=table_field_name_str+","
			 }
			 
			 path_attribs="//dia:attribute[@name=\"unique\"]/dia:boolean[@val=\"true\"]"
			 table_field_is_unique,_:=xmlquery.Query(table_field_as_xml_node,path_attribs)
			 if(table_field_is_unique!=nil) {
				 unique_field_names_list_str+=table_field_name_str+","
			 }
		 }
		 
		 table_sql_query=strings.Trim(table_sql_query,",") + ");"
		 fmt.Println(table_sql_query)
		 if (len(primary_keys_list_str)>0) {
			 primary_keys_query:="alter table "+table_name_str+" add constraints PK_"+table_name_str+" primary key ("+strings.Trim(primary_keys_list_str,",")+");"
			 primary_keys_list_str=""
			 fmt.Println(primary_keys_query)
			 primary_keys_query=""
		 }
		 
		 if (len(unique_field_names_list_str)>0) {
			 unique_query:="alter table "+table_name_str+" add constraints UNIQUE_"+table_name_str+" unique ("+strings.Trim(unique_field_names_list_str,",")+");"
			 unique_field_names_list_str=""
			 fmt.Println(unique_query)
			 unique_query=""
		 }
		 
		 
		 
		 table_sql_query=""
		 
    }

}
