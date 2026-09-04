program helloworld;
var
   x, y : string;

procedure foo(x : string);
   begin
      writeln(x);
   end;

begin
   x := 'Hello';
   y := 'World';

   foo(x + ', ' + y + '!')
end.
