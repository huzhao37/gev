namespace go cluster

 struct Result {
     1: i64 div;
     2: i64 mod;
 }

 service DivMod {
     Result DoDivMod(1:i64 arg1, 2:i64 arg2);
 }

  service DivMod2 {
      Result DoDivMod2(1:i64 arg1, 2:i64 arg2);
  }