PROGRAM HelloWorld;
VAR
   a, b : INTEGER;
   c    : REAL;
   x, y : STRING;

PROCEDURE Print(x : String);
   BEGIN
      WriteLn(x)
   END;

BEGIN
   a := 5;
   b := a;
   c := 1.2;

   x := 'Hello';
   y := 'World';

   Print(x + ', ' + y + '!')
END.
