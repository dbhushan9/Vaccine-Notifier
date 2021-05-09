*Vaccine Centers \(45\+ Yrs\)*

```
{{range .CenterDetails.Age45}}
1. Name     : {{.Name}}
2. Capacity : {{.AvailableCapacity}}
3. Vaccine  : {{.Vaccine}}
4. Fee Type : {{.FeeType}}
---------------------------------------------
{{end}}
```

*Vaccine Centers \(18\-44 Yrs\)*

```
{{range .CenterDetails.Age18}}
1. Name     : {{.Name}}
2. Capacity : {{.AvailableCapacity}}
3. Vaccine  : {{.Vaccine}}
4. Fee Type : {{.FeeType}}
---------------------------------------------
{{end}}
```