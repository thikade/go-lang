<!doctype html> 
<htmllang="en"> 
<head> 
   <meta charset="UTF-8"> 
   <title>Using Go HTML Templates</title> 
   <style> 
         html { 
               font-size: 16px; 
         } 
         table, th, td { 
         border: 3px solid red; 
         } 
   </style> 
</head> 
<body> 
 
<h2 alight="center">Presenting Dynamic content read from filesystem!</h2> 
 
<table> 
   <thead> 
         <tr> 
               <th>Animal</th> 
               <th>Age</th> 
               <th>?</th> 
         </tr> 
   </thead> 
   <tbody> 
{{ range . }} 
<tr> 
   <td> {{ .Animal }} </td> 
   <td> {{ .Age }} </td> 
   <td> {{ .Age }} </td> 
</tr> 
{{ end }} 
   </tbody> 
</table> 
 
</body> 
</html> 