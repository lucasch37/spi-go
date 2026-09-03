PROGRAM HelloWorld;
VAR
   a, b : INTEGER;
   c    : REAL;
   x, y : STRING;

BEGIN
   a := 5;
   b := a;
   c := 1.2;

   x := 'Hello';
   y := 'World';

   writeln(x + ', ' + y + '!')
END.
