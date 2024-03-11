# SQLit

## Online demo
Frontend repo here: https://github.com/jnafolayan/sqlit-web

Deployment URL: http://sqlit.vercel.app 

## Potential improvements
- Currently, I'm reusing the `ast` objects in the evaluator. This creates unwanted artifacts when describing evaluated objects. An improvement would be to create an object ontology for the evaluator.
- Abstract as many type-specific operations into a package. 
- Refactor the parser to more explicitly define its components.
