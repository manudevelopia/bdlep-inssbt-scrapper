# bdlep-inssbt-scrapper

This tool scraps http://bdlep.inssbt.es/LEP/ just to retrieve chemical composes information (NºCE, NºCAS, VLA-Xx, Notes, Warning advices, ..) and provide it to 
Prevengic https://prevengic.herokuapp.com. 


### why?
As far I know, LEP Chemical information is enclosed in http://bdlep.inssbt.es/LEP/ website and PDF http://www.inssbt.es/InshtWeb/Contenidos/Documentacion/LEP%20_VALORES%20LIMITE/Valores%20limite/Limites2018/Limites2018.pdf
Those formats are nor very useful I you pretend to use them intensively.

### Golang?
Why not!! By using it makes me to get outside from my daily work language.

### Process
- Read main composes index at http://bdlep.inssbt.es/LEP/vlaallpr.jsp?Bloque=1&submit=Listado+completo+Agentes+Qu%C3%ADmicos
    - Extract: compose name and link to information page

- Read each compose information page http://bdlep.inssbt.es/LEP/vlapr.jsp?ID=1&nombre=Aceite%20mineral%20refinado,%20nieblas
    - Extract: NºCAS, NºCE, VLA-EC/ED, Notes, Warning Advices
    
- Store information by each compose

- Report compose information Prevengic application

### Result
Gathered data will be added to the current version of Prevengic application. You can find last version of code at https://github.com/manudevelopia/prevengic and use application at https://prevengic.herokuapp.com
