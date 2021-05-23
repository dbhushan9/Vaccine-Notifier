{{if .CenterDetails.Age45}}
_Vaccine Centers \(45\+ Yrs\)_

```
{{range .CenterDetails.Age45}}
1. Name     : {{.Name}}
2. Pincode : {{.Pincode}}
3. Capacity : {{.AvailableCapacity}}
4. Vaccine  : {{.Vaccine}}
5. Fee Type : {{.FeeType}}
---------------------------------------------
{{end}}
```
{{end}}
{{if .CenterDetails.Age45}}
_Vaccine Centers \(18\-44 Yrs\)_

```
{{range .CenterDetails.Age18}}
1. Name     : {{.Name}}
2. Pincode : {{.Pincode}}
3. Capacity : {{.AvailableCapacity}}
4. Vaccine  : {{.Vaccine}}
5. Fee Type : {{.FeeType}}
---------------------------------------------
{{end}}
```
{{end}}