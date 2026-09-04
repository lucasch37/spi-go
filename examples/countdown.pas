program countdown;

procedure countdown(num : integer);
   begin
      if num <= 0 then
         writeln('Go!')
      else
         begin
            writeln(num);
            countdown(num-1)
         end
   end;

begin
   countdown(10)
end.
