package entity

type SceneOpts struct {
	ID     string `msgpack:"id"`
	Name   string `msgpack:"name"`
	MinX   int    `msgpack:"minX"`
	MaxX   int    `msgpack:"maxX"`
	ContsX int    `msgpack:"contsX"`
	MinY   int    `msgpack:"minY"`
	MaxY   int    `msgpack:"maxY"`
	ContsY int    `msgpack:"contsY"`
}
