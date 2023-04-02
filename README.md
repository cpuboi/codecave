# CodeCave #

Searches a file for a space of unused bytes and inserts a "hidden" message.  

## Insertion ##
When data is sent to CodeCave it looks for a cave of zeroes of sufficient size, inserts the data beginning and ending with a special pattern.  
The special pattern is based on the first 8 bytes of an MD5 sum that in turn is based on the first 48 bytes of the file.  


## Extraction ##
When the special pattern is found, CodeCave keeps looking until the same pattern is found again (within a set length of bytes)  
The data is extracted and sent to STDOUT  


## Usage ##
```
Usage of ./codecave:
  -a	Analyze file, show all codecaves in csv format, combine with -f
  -d	Decode message from file
  -e string
    	Encode message into file, provide message
  -f string
    	Input file (default "./example_input_file.bin")
  -o string
    	Output file (default "./example_output_file.bin")
  -t	Test run, dont actually write output, needs to be combined with -e
  -v	Verbose output
```
## Examples ##
* **Encode data**  
``codecave -e "hide this message" -f ./image.gif -o ./image2.gif``  


* **Decode data**  
``codecave -d -f ./image2.gif``


* **Find files with code caves**  
``find ./directory -type f  -size -10M -exec ./codecave -f {} -a \; | grep -v "0,0,0"``  
This example searches for files with less than 10 megabyte size and removes all files without code caves (the grep -v "0,0,0" part)    
The analyze function returns: ``filename,size,start,end`` if size,start and end are 0 then no caves were found.  

