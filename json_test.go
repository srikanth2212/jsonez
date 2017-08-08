package jsonez

import (
	"fmt"
	"strings"
	"testing"
)

func TestSimpleParse(t *testing.T) {
	input := []byte(`{
		"outer":	{
				"val1":	"foo",
				"val2":	"bar",
				"val3":	1234,
				"val4":	225.1245,
				"val5":	[
							1,
							2,
							3,
							4,
							5
				]
			}
		}`)

	g, err := GoJSONParse(input)

	if err != nil {
		t.Errorf("%s: GoJSONParse failed with error %s", funcName(), err)
		return
	}

	output := GoJSONPrint(g)
	if output == nil {
		t.Errorf("%s: Print failed", funcName())
	}

	o, err := g.Get("outer", "val1")

	if err != nil {
		t.Errorf("%s: key val1 not found", funcName())
	} else if strings.Compare(o.Valstr, "foo") != 0 {
		t.Errorf("%s: key val1 has incorrect value %s", funcName(), o.Valstr)
	} else {
		t.Logf("%s: key val1 has the correct value %s", funcName(), o.Valstr)
	}

	_, err = g.GetIntVal("outer", "val1")
	if err != nil {
		t.Logf("%s: GetIntVal for key val1 failed as expected with error %s", funcName(), err)
	} else {
		t.Errorf("%s: GetIntVal for key val1 didn't fail as expected", funcName())
	}

	_, err = g.GetDoubleVal("outer", "val1")
	if err != nil {
		t.Logf("%s: GetIntVal for key val1 failed as expected with error %s", funcName(), err)
	} else {
		t.Errorf("%s: GetIntVal for key val1 didn't fail as expected", funcName())
	}

	i, err := g.GetUIntVal("outer", "val3")
	if err != nil {
		t.Errorf("%s: GetUIntVal for key val3 failed with error %s", funcName(), err)
	} else if i != 1234 {
		t.Errorf("%s: GetUIntVal for key val3 returned %d while expected value was 1234", funcName(), i)
	} else {
		t.Logf("%s: GetUIntVal for key val3 returned 1234 as expected", funcName())

	}

	d, err := g.GetDoubleVal("outer", "val4")
	if err != nil {
		t.Errorf("%s: GetDoubleVal for key val4 failed with error %s", funcName(), err)
	} else if d != 225.1245 {
		t.Errorf("%s: GetDoubleVal for key val4 returned %f while expected value "+
			"was 225.1245", funcName(), d)
	} else {
		t.Logf("%s: GetDoubleVal for key val4 returned 225.1245 as expected",
			funcName())
	}

	a, err := g.Get("outer", "val5")
	if a.GetArraySize() != 5 {
		t.Errorf("%s: Arraysize of the object with key val5 returned %d while "+
			"expected value was 5", funcName(), a.GetArraySize())
	} else {
		t.Logf("%s: Arraysize of the object with key val5 returned 5 as expected",
			funcName())
	}

	return
}

func TestBuilder(t *testing.T) {
	input := []byte(`{
"outer":  {
"val1": "foo",
"val2": "bar",
"val3": 1234,
"val4": 225.1245,
"val5": [
1,
2,
3,
4,
5
]
}
}`)

	g, err := GoJSONParse(input)

	if err != nil {
		t.Errorf("%s: GoJSONParse failed with error %s", funcName(), err)
		return
	} else {
		t.Logf("%s: GoJSONParse completed", funcName())
	}

	err = g.AddVal(100, "outer", "val6")
	if err != nil {
		t.Errorf("%s: AddVal failed with error %s", funcName(), err)
	}

	v, err := g.GetUIntVal("outer", "val6")
	if err != nil {
		t.Errorf("%s: Searching for key val6 within outer failed with error %s", funcName(), err)
	} else if v != 100 {
		t.Errorf("%s: Val6 returned %d while expected was 100", funcName(), v)
	} else {
		t.Logf("%s: key val6 was found within outer and returned 100", funcName())
	}

	fmt.Println("printing after adding val6")
	fmt.Println(string(GoJSONPrint(g)))

	err = g.AddVal(245.67, "outer", "val7")
	if err != nil {
		t.Errorf("%s: AddVal failed with error %s", funcName(), err)
	}

	vd, err := g.GetDoubleVal("outer", "val7")
	if err != nil {
		t.Errorf("%s: Searching for key val7 within outer failed with error %s",
			funcName(), err)
	} else if vd != 245.67 {
		t.Errorf("%s: Val6 returned %f while expected was 100", funcName(), vd)
	} else {
		t.Logf("%s: key val6 was found within outer and returned 100", funcName())
	}

	err = g.AddVal("hello world", "outer", "val8")
	if err != nil {
		t.Errorf("%s: AddVal failed with error %s", funcName(), err)
	}

	_, err = g.Get("outer", "val8")
	if err != nil {
		t.Errorf("%s: Searching for key val8 within outer failed with error %s", funcName(), err)
	} else {
		t.Logf("%s: key val8 found within outer", funcName())
	}

	err = g.AddVal(true, "outer", "val9")
	if err != nil {
		t.Errorf("%s: AddVal failed with error %s", funcName(), err)
	}

	_, err = g.Get("outer", "val9")
	if err != nil {
		t.Errorf("%s: Searching for key val8 within outer failed with error %s", funcName(), err)
	} else {
		t.Logf("%s: key val9 found within outer", funcName())
	}

	err = g.AddToArray(100, "outer", "val10")
	if err != nil {
		t.Errorf("%s: AddToArray failed with error %s", funcName(), err)
	}

	_, err = g.Get("outer", "val10")
	if err != nil {
		t.Errorf("%s: Searching for key val10 within outer failed with error %s",
			funcName(), err)
	} else {
		t.Logf("%s: key val10 found within outer", funcName())
	}
	fmt.Println("printing after adding val10")
	fmt.Println(string(GoJSONPrint(g)))

	err = g.AddToArray(200.25, "outer", "val10")
	if err != nil {
		t.Errorf("%s: AddToArray failed with error %s", funcName(), err)
	}

	_, err = g.Get("outer", "val10")
	if err != nil {
		t.Errorf("%s: Searching for key val10 within outer failed with error %s",
			funcName(), err)
	} else {
		t.Logf("%s: key val10 found within outer", funcName())
	}

	err = g.AddToArray("hello world", "outer", "val10")
	if err != nil {
		t.Errorf("%s: AddToArray failed with error %s", funcName(), err)
	}

	_, err = g.Get("outer", "val10")
	if err != nil {
		t.Errorf("%s: Searching for key val10 within outer failed with error %s",
			funcName(), err)
	} else {
		t.Logf("%s: key val10 found within outer", funcName())
	}

	err = g.AddToArray(true, "outer", "val10")
	if err != nil {
		t.Errorf("%s: AddToArray failed with error %s", funcName(), err)
	}

	arr, err := g.Get("outer", "val10")
	if err != nil {
		t.Errorf("%s: Searching for key val10 within outer failed with error %s",
			funcName(), err)
	} else {
		t.Logf("%s: key val10 found within outer", funcName())
	}

	entry, err := arr.GetArrayElemByIndex(1)

	if err != nil {
		t.Errorf("%s: Searching for array entry at location 1 failed with error %s", funcName(), err)
	} else if entry.Valdouble != 200.25 {
		t.Errorf("%s: Array entry at location is %f while expeced is 200.25", funcName(), entry.Valdouble)
	} else {
		t.Logf("%s Value 200.25 found at location 1", funcName())
	}

	fmt.Println(string(GoJSONPrint(g)))

	err = arr.DelIndexFromArray(3)
	if err != nil {
		t.Errorf("%s: Deleting element at location 2 failed with error %s", funcName(), err)
	} else {
		t.Logf("%s: Deleted element at location 2", funcName())
	}

	fmt.Println(string(GoJSONPrint(g)))

	err = arr.DelIndexFromArray(10)
	if err != nil {
		t.Logf("%s: Deleting element at location 10 failed as expected with error %s", funcName(), err)
	} else {
		t.Errorf("%s: Deleting element at location 10 didn't fail as expected", funcName())
	}

	_, err = g.Get("outer")
	if err != nil {
		t.Errorf("%s: Error fetching object outer", funcName())
	}

	fmt.Println(string(GoJSONPrint(g)))

	err = g.DelFromArray(3, "outer", "val5")
	if err != nil {
		t.Errorf("%s: Deleting entry 3 within val5 failed with error %s", funcName(), err)
	}

	err = g.DelVal("outer", "val10")

	err = g.AddVal(-15, "outer", "val11")

	fmt.Println(string(GoJSONPrint(g)))

	v1, err := g.GetIntVal("outer", "val11")
	if err != nil {
		t.Errorf("%s: Fetching val11 failed with error %s", funcName(), err)
	} else if v1 != -15 {
		t.Errorf("%s: val11 doesn't have -15 as value", funcName())
	} else {
		t.Logf("%s: val11 has -15 as expected", funcName())
	}
	/*
	   if err != nil {
	   t.Logf("%s: Deleting child object failed as expected with error %s", funcName(), err)
	   } else {
	   t.Errorf("%s: Deleting child object did not fail as expected", funcName())
	   }
	*/

}

func TestComplexParse(t *testing.T) {
	input := `[
  {
    "_id": "58fab57c4d4a638a74d8962c",
    "index": 0,
    "guid": "291dc9a5-ee3f-4cbb-aae4-e9210bfef8fe",
    "isActive": false,
    "balance": "$3,672.11",
    "picture": "http://placehold.it/32x32",
    "age": 37,
    "eyeColor": "green",
    "name": {
      "first": "Graciela",
      "last": "Britt"
    },
    "company": "ZYTREX",
    "email": "graciela.britt@zytrex.io",
    "phone": "+1 (846) 578-3461",
    "address": "547 Regent Place, Chamizal, Mississippi, 5971",
    "about": "Ipsum Lorem ea pariatur enim. Tempor excepteur amet sit excepteur. Dolore ut duis irure ad laborum ex ut laboris enim tempor sunt duis quis dolor. Labore mollit excepteur culpa dolore. Excepteur sint laborum voluptate ex eiusmod veniam proident. Tempor deserunt laborum ea id Lorem dolor in labore eu laboris id irure consequat occaecat.",
    "registered": "Monday, March 23, 2015 8:26 AM",
    "latitude": "80.365295",
    "longitude": "176.046508",
    "tags": [
      "in",
      "in",
      "deserunt",
      "est",
      "cillum"
    ],
    "range": [
      0,
      1,
      2,
      3,
      4,
      5,
      6,
      7,
      8,
      9
    ],
    "friends": [
      {
        "id": 0,
        "name": "Everett Levine"
      },
      {
        "id": 1,
        "name": "Noble Chang"
      },
      {
        "id": 2,
        "name": "Knox Combs"
      }
    ],
    "greeting": "Hello, Graciela! You have 8 unread messages.",
    "favoriteFruit": "banana"
  },
  {
    "_id": "58fab57d46ea2075188a0319",
    "index": 1,
    "guid": "d47a9077-6f2a-4c19-ab65-4cbd6561e199",
    "isActive": true,
    "balance": "$1,291.86",
    "picture": "http://placehold.it/32x32",
    "age": 35,
    "eyeColor": "green",
    "name": {
      "first": "Avila",
      "last": "Becker"
    },
    "company": "DIGINETIC",
    "email": "avila.becker@diginetic.biz",
    "phone": "+1 (862) 484-2144",
    "address": "810 Bedford Avenue, Dellview, Marshall Islands, 5054",
    "about": "Ut laboris consequat nulla consectetur id laborum fugiat. Quis aliqua sint dolore ea ex labore aliqua ea aliquip elit. Proident dolor commodo voluptate excepteur tempor. Ipsum non aliqua esse cillum amet veniam incididunt minim adipisicing. Fugiat non occaecat velit adipisicing exercitation amet do aute pariatur proident ullamco exercitation. Culpa veniam eiusmod cillum do sunt tempor culpa cupidatat aliqua exercitation sint veniam eiusmod.",
    "registered": "Friday, September 4, 2015 7:46 AM",
    "latitude": "-44.466674",
    "longitude": "-6.10677",
    "tags": [
      "incididunt",
      "sunt",
      "commodo",
      "cillum",
      "nisi"
    ],
    "range": [
      0,
      1,
      2,
      3,
      4,
      5,
      6,
      7,
      8,
      9
    ],
    "friends": [
      {
        "id": 0,
        "name": "Erica Pearson"
      },
      {
        "id": 1,
        "name": "Sherman Lawrence"
      },
      {
        "id": 2,
        "name": "Compton Conrad"
      }
    ],
    "greeting": "Hello, Avila! You have 6 unread messages.",
    "favoriteFruit": "apple"
  },
  {
    "_id": "58fab57d6b128b1745c77694",
    "index": 2,
    "guid": "2e03abb5-bffa-4635-9e90-d5a76ee9ed4d",
    "isActive": false,
    "balance": "$1,262.46",
    "picture": "http://placehold.it/32x32",
    "age": 33,
    "eyeColor": "green",
    "name": {
      "first": "Queen",
      "last": "Barrett"
    },
    "company": "XOGGLE",
    "email": "queen.barrett@xoggle.ca",
    "phone": "+1 (876) 588-3641",
    "address": "685 Coyle Street, Caberfae, West Virginia, 8475",
    "about": "Laborum culpa ut fugiat esse velit consequat. Incididunt id aliquip reprehenderit commodo eu veniam. Ut dolor esse occaecat sunt cillum.",
    "registered": "Wednesday, October 12, 2016 1:09 AM",
    "latitude": "-87.386704",
    "longitude": "-37.456147",
    "tags": [
      "cupidatat",
      "elit",
      "ullamco",
      "eiusmod",
      "esse"
    ],
    "range": [
      0,
      1,
      2,
      3,
      4,
      5,
      6,
      7,
      8,
      9
    ],
    "friends": [
      {
        "id": 0,
        "name": "Sheree Wilson"
      },
      {
        "id": 1,
        "name": "Ana Huff"
      },
      {
        "id": 2,
        "name": "Anna Larson"
      }
    ],
    "greeting": "Hello, Queen! You have 7 unread messages.",
    "favoriteFruit": "strawberry"
  },
  {
    "_id": "58fab57d8737b95d416dc5cc",
    "index": 3,
    "guid": "5bccf39c-c048-4c65-9fc1-0656c54ce723",
    "isActive": true,
    "balance": "$3,597.54",
    "picture": "http://placehold.it/32x32",
    "age": 32,
    "eyeColor": "blue",
    "name": {
      "first": "Jolene",
      "last": "Farrell"
    },
    "company": "EMOLTRA",
    "email": "jolene.farrell@emoltra.me",
    "phone": "+1 (971) 522-3511",
    "address": "764 Sedgwick Street, Carbonville, Vermont, 9876",
    "about": "Nostrud aliquip sint cillum sint mollit ullamco laborum consectetur laborum nisi ut consectetur nisi incididunt. Laboris nostrud ad elit proident cupidatat do veniam. Consequat adipisicing sint aliquip adipisicing ullamco cupidatat sint deserunt ex laborum esse voluptate sint enim. Incididunt magna excepteur labore do eu velit nostrud aliqua nulla. Qui veniam exercitation commodo occaecat dolor adipisicing exercitation elit. Nulla nisi aliquip laborum cupidatat cupidatat.",
    "registered": "Thursday, January 5, 2017 2:10 AM",
    "latitude": "-23.573173",
    "longitude": "39.236471",
    "tags": [
      "velit",
      "do",
      "nulla",
      "cupidatat",
      "proident"
    ],
    "range": [
      0,
      1,
      2,
      3,
      4,
      5,
      6,
      7,
      8,
      9
    ],
    "friends": [
      {
        "id": 0,
        "name": "Lisa Montoya"
      },
      {
        "id": 1,
        "name": "Thomas Hardin"
      },
      {
        "id": 2,
        "name": "Leonor Adams"
      }
    ],
    "greeting": "Hello, Jolene! You have 10 unread messages.",
    "favoriteFruit": "strawberry"
  },
  {
    "_id": "58fab57d812649a354bbc5d0",
    "index": 4,
    "guid": "36370e91-ed8c-4754-9c2c-508071731846",
    "isActive": false,
    "balance": "$2,370.42",
    "picture": "http://placehold.it/32x32",
    "age": 25,
    "eyeColor": "green",
    "name": {
      "first": "Bartlett",
      "last": "Beasley"
    },
    "company": "MENBRAIN",
    "email": "bartlett.beasley@menbrain.us",
    "phone": "+1 (999) 470-3129",
    "address": "856 Clara Street, Gerton, Nebraska, 7588",
    "about": "Nostrud magna dolore minim eu voluptate. Do aliquip velit aliqua dolor do ipsum ad tempor consequat sint irure laborum Lorem voluptate. Officia cillum sunt occaecat duis eu ex adipisicing et.",
    "registered": "Monday, August 25, 2014 3:41 AM",
    "latitude": "74.302983",
    "longitude": "-72.936819",
    "tags": [
      "incididunt",
      "cillum",
      "enim",
      "enim",
      "aute"
    ],
    "range": [
      0,
      1,
      2,
      3,
      4,
      5,
      6,
      7,
      8,
      9
    ],
    "friends": [
      {
        "id": 0,
        "name": "Mays Boyle"
      },
      {
        "id": 1,
        "name": "Vang Copeland"
      },
      {
        "id": 2,
        "name": "Ruby Camacho"
      }
    ],
    "greeting": "Hello, Bartlett! You have 9 unread messages.",
    "favoriteFruit": "apple"
  },
  {
    "_id": "58fab57d13362fc7a9bb353e",
    "index": 5,
    "guid": "81b97d09-cf01-4ec5-bfe9-399f8b1a4b61",
    "isActive": false,
    "balance": "$1,624.29",
    "picture": "http://placehold.it/32x32",
    "age": 29,
    "eyeColor": "brown",
    "name": {
      "first": "Vilma",
      "last": "Parsons"
    },
    "company": "PROVIDCO",
    "email": "vilma.parsons@providco.org",
    "phone": "+1 (842) 558-2965",
    "address": "606 Schenck Avenue, Crenshaw, Guam, 4135",
    "about": "Quis dolor ipsum laborum laborum mollit. Velit nisi exercitation nostrud reprehenderit deserunt labore incididunt culpa elit aliquip veniam adipisicing amet occaecat. Laborum ipsum minim consequat consectetur ex labore deserunt. Ea eiusmod pariatur cupidatat sint qui exercitation culpa officia deserunt laboris do mollit elit labore. Consectetur exercitation dolore amet duis minim non tempor. Laborum do eiusmod id anim incididunt adipisicing labore incididunt. Irure deserunt qui aliquip Lorem.",
    "registered": "Saturday, August 29, 2015 4:17 PM",
    "latitude": "-87.299616",
    "longitude": "-99.445812",
    "tags": [
      "mollit",
      "non",
      "pariatur",
      "laborum",
      "ad"
    ],
    "range": [
      0,
      1,
      2,
      3,
      4,
      5,
      6,
      7,
      8,
      9
    ],
    "friends": [
      {
        "id": 0,
        "name": "Bernadette Jackson"
      },
      {
        "id": 1,
        "name": "Medina Leonard"
      },
      {
        "id": 2,
        "name": "Eloise Melton"
      }
    ],
    "greeting": "Hello, Vilma! You have 8 unread messages.",
    "favoriteFruit": "apple"
  },
  {
    "_id": "58fab57d202be322c171d9d5",
    "index": 6,
    "guid": "aef23887-74b3-4ee1-9c2f-dcaf6af1c39e",
    "isActive": true,
    "balance": "$3,852.70",
    "picture": "http://placehold.it/32x32",
    "age": 39,
    "eyeColor": "blue",
    "name": {
      "first": "Christa",
      "last": "Franco"
    },
    "company": "MAGNEATO",
    "email": "christa.franco@magneato.tv",
    "phone": "+1 (910) 570-2105",
    "address": "805 Eckford Street, Fingerville, Texas, 9366",
    "about": "Ad est velit est proident ex. Officia enim aliqua adipisicing nostrud incididunt culpa incididunt. Reprehenderit sit proident magna reprehenderit laboris aliquip qui amet irure cillum. Sit sunt occaecat ad ea nisi nostrud id aute aliqua. Deserunt sunt aliqua occaecat culpa mollit labore commodo sunt do. Dolore ea enim incididunt et fugiat quis sint duis duis. Elit nostrud mollit mollit nulla qui veniam nostrud.",
    "registered": "Friday, August 8, 2014 7:01 PM",
    "latitude": "-17.24687",
    "longitude": "13.312776",
    "tags": [
      "cupidatat",
      "et",
      "tempor",
      "laboris",
      "consequat"
    ],
    "range": [
      0,
      1,
      2,
      3,
      4,
      5,
      6,
      7,
      8,
      9
    ],
    "friends": [
      {
        "id": 0,
        "name": "Owens Wong"
      },
      {
        "id": 1,
        "name": "Poole Gordon"
      },
      {
        "id": 2,
        "name": "Fannie Maynard"
      }
    ],
    "greeting": "Hello, Christa! You have 10 unread messages.",
    "favoriteFruit": "strawberry"
  }
]
`

	g, err := GoJSONParse([]byte(input))

	if err != nil {
		t.Errorf("%s: GoJSONParse failed with error %s", funcName(), err)
		return
	}

	output := GoJSONPrint(g)

	if output == nil {
		t.Errorf("%s: GoJSONPrint failed as output was empty", funcName())
		return
	}

	return
}
