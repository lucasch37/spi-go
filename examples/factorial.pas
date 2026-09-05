program FactorialProgram;

var
  num: Integer;

function Factorial(n: Integer): Integer;
begin
  if n = 0 then
    Factorial := 1
  else
    Factorial := n * Factorial(n - 1);
end;

procedure PrintFactorials(n : Integer);
begin
  writeln(n, '! = ', Factorial(n));
  if n < 10 then
    PrintFactorials(n + 1);
end;

begin
  num := 1;

  PrintFactorials(num);
end.
