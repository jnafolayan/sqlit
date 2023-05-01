# SQLit

## Online demo
http://sqlit.vercel.app

## Potential improvements
- Currently I'm reusing the `ast` objects in the evaluator. This creates unwanted artifacts when describing evaluated objects. An improvement would be to create an object ontology for the evaluator.
- Abstract as many type-specific operations into a package. 
- Refactor the parser more explicitly define its components.