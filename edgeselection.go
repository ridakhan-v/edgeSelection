package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	//"github.com/paulmach/orb"
)


type MultiPolygon struct {
	coordinates [][][][]float64
}

type Geojson struct {
	Type        string          `json:"type"`
	Coordinates json.RawMessage `json:"coordinates"`
	//Point       Point           `json:"-"`
	//Line        Line            `json:"-"`
	MultiPolygon     MultiPolygon         `json:"-"`
}



func main() {
	//fileName := "test3.json"
	//fileName :="ec01b2d6-bdd4-11ea-95d2-005056899273.json"
	//fileName := "8b6a5212-a447-11ea-b6da-005056899273.json"
	//fileName :="94d71cfe-8f9d-48b1-8f79-1807ecefebee.json"
	//fileName := "c83907d3-13eb-484c-921d-f6afa4e13721.json"
	fileName := "315143ef-ef19-48ca-8b0d-59f68cbe22ac.json"
	//fileName :="e221df34-3ae3-11eb-a648-005056899273.json"

	jsonFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully opened shape.json")

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	var poly Geojson
	err2:= poly.UnmarshalJSON([]byte(byteValue))
	if err2 != nil {
		fmt.Println(err2)
	}
	
	P2 := poly.MultiPolygon.coordinates[0][0]
	
	//var P2 [][2]float64
	//P2 = [6][2]float64{  {20, 40}, {40, 40}, {35.22282608695652328 ,26.14130434782609314}, {30, 10},{10, 20},{20, 40}}

	//P2 :=[13][2] float64 {{83.16664417613638705, 43.07392666903406564}, {96.9884623579545746, 35.74301757812497016},
	//{104.7011896306818528 ,19.24847212357951065}, {103.86118963068184939, 5.73210848721587496}, {92.48300781250003411, 0.15756303267042426},
	//{74.53755326704548168, -1.44607333096594459}, {52.46846235795456437, -0.60607333096594118}, {42.99937144886365559, 3.51756303267042369},
	//{35.43937144886366042, 2.06665394176133077}, {45.21391690340911396, 13.97938121448860471}, {47.19937144886365843, 27.03756303267042682},
	//{53.00300781250002302,41.62301757812497272}, {83.16664417613638705, 43.07392666903406564}}

	//P2 := [17][4]float64{{83.3858695652174049, 80.33695652173912549}, {92.3858695652174049, 78.64130434782607892}, {97.34239130434784215, 71.59782608695651618},
	//{96.55978260869565588, 61.55434782608694633}, {93.95108695652174902, 52.1630434782608603}, {85.99456521739131176, 48.64130434782607892},
	//{77.3858695652174049, 47.20652173913042304}, {71.25543478260870245, 48.51086956521738358}, {64.34239130434782794, 51.51086956521738358},
	//{60.42934782608696764, 56.0760869565217277}, {59.51630434782609314, 61.03260869565216495},{ 60.42934782608696764, 67.55434782608693922},
	//{60.82065217391304657, 72.64130434782607892}, {64.08152173913043725, 78.77173913043476716}, {70.2119565217391397, 80.46739130434781373},
	//{74.64673913043479558, 81.11956521739129755}, {83.3858695652174049, 80.33695652173912549} }

	var set [][]float64
	set = append(set, P2[0], P2[1])

	currentSetLength := length(set[0], set[1])
	longestLength := 0.0
	var longestSet [][]float64
	var threshold float64
	threshold = 220


	for i := 2; i < len(P2); i++ {

		t := angleBetweenVectors(P2[i], P2[i-1], P2[i-2])
		l := length(P2[i-1], P2[i])

		if l < 10 {
			P2 = RemoveIndex(P2, i)
			i -= 1
			continue}

		fmt.Println("angele:",t, "currentSetLength=", currentSetLength, "longestLength=", longestLength, "i", i,  P2[i-1], P2[i])
		fmt.Println(length(P2[i-1], P2[i]), "length")


		if i == len(P2)-1{
			if t <= threshold {
				currentSetLength += length(set[len(set)-1], P2[i])
				set = append(set, P2[i])

				if currentSetLength > longestLength {
					longestLength = currentSetLength
					longestSet = set
				}
				continue
			}

			if t > threshold {
				if currentSetLength > longestLength {
					longestLength = currentSetLength
					longestSet = set
					set = nil
					set = append(set, P2[i-1], P2[i])
					currentSetLength = length(set[0], set[1])
					fmt.Println(currentSetLength, "current length", set, "current set")
				}
				if currentSetLength <= longestLength {
					set = nil
					set = append(set, P2[i-1], P2[i])
					currentSetLength = length(set[0], set[1])
				}

			}
		continue

		}

		if t <= threshold {
			currentSetLength += length(set[len(set)-1], P2[i])
			set = append(set, P2[i])
			fmt.Println("**")
			continue
		}

		if t > threshold {
			if currentSetLength > longestLength {
				longestLength = currentSetLength
				longestSet = set
				set = nil
				set = append(set, P2[i-1], P2[i])
				currentSetLength = length(set[0], set[1])
				fmt.Println("-----------------current length=",currentSetLength, "current set=", set, "-------Longest set", longestSet, "----------")
				continue
			}
			if currentSetLength <= longestLength {
				set = nil
				set = append(set, P2[i-1], P2[i])
				currentSetLength = length(set[0], set[1])
			}
			continue
		}
	}
	fmt.Println("current Set Length", currentSetLength)
	fmt.Println("longestLength", longestLength)
	fmt.Println("Longest set", longestSet)
	fmt.Println("Longest set size", len(longestSet))

}

// Haversine formula to calculate great-circle distance - the shortest distance over
// the earth's surface - in meters
func length(P1, P2 []float64) float64 {
	DLat := (P2[1]-P1[1])*math.Pi/180
	DLon := (P2[0]-P1[0])*math.Pi/180
	a := math.Pow(math.Sin(DLat/2),2)  + (math.Cos(P1[1]*math.Pi/180)*math.Cos(P2[1]*math.Pi/180)*math.Pow(math.Sin(DLon/2),2))
	c := 2 * math.Atan2(math.Sqrt(a),math.Sqrt(1-a))
	d :=  6371000 * c // distance in meters
	return d
}

func angleBetweenVectors(a, b, c []float64) float64 {
	BAx := a[0] - b[0]
	BAy := a[1] - b[1]

	BCx := c[0] - b[0]
	BCy := c[1] - b[1]

	BAmag := math.Sqrt(BAx*BAx + BAy*BAy)
	BCmag := math.Sqrt(BCx*BCx + BCy*BCy)

	angle := math.Acos((BAx*BCx + BAy*BCy) / (BAmag * BCmag))
	return 360 - (angle * 180 / math.Pi)
}


func (g *Geojson) UnmarshalJSON(b []byte) error {
	type Alias Geojson
	aux := (*Alias)(g)

	err := json.Unmarshal(b, &aux)

	if err != nil {
		return err
	}

	switch g.Type {
	//case "Point":
	//	err = json.Unmarshal(g.Coordinates, &g.Point.Coordinates)
	//case "LineString":
	//	err = json.Unmarshal(g.Coordinates, &g.Line.Points)
	case "MultiPolygon":
		err = json.Unmarshal(g.Coordinates, &g.MultiPolygon.coordinates)

	}
	g.Coordinates = []byte(nil)
	return err
}

func RemoveIndex(s [][]float64, index int) [][]float64 {
	return append(s[:index], s[index+1:]...)
}