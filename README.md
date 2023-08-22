# dia2sql-golang

This is a Go language program that opens a Dia diagram file that contains Database shapes and generates the SQL statements of "create table.." from each shape that represents a table.

# Building the program
1. Make sure Go is installed on your system.
1. Change in your terminal to the diectory of the project, then type:
`go build dia2sql`

# Executing the program
type :
`./dia2sql your_dia_diagram.dia`

if you type: `./dia2sql -v your_dia_diagram.dia` it will show Dia file analysis

the file _school.dia_ is a sample file to use as a test
