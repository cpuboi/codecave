# CodeCave #

Finds a code cave of zeroes   
Inserts data in to code caves  
Extracts data from code caves  


## Insertion ##
When data is sent to CodeCave it looks for a cave of sufficient size, inserts the data beginning and ending with a special pattern.  
It then saves the file.  
The special pattern is based on the first 8 bytes of an MD5 sum that in turn is based on the first 48 bytes of the file.  


## Extraction ##

When the special pattern is found, CodeCave keeps looking until the same pattern is found again (within a set length of bytes)  
The data inbetween is extracted and sent to STDOUT  

