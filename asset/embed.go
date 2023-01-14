package asset

import (
    _ "embed"
)

// Sounds
var (
    //go:embed fire.wav
    Fire_wav []byte

    //go:embed jab.wav
    Hit_wav []byte

    //go:embed menu.wav
    Menu_wav []byte

    //go:embed ship_destroy.wav
    Destroyed_wav []byte

    //go:embed thruster.wav
    Thruster_wav []byte

    //go:embed wave.wav
    Wave_wav []byte
)

// Sprites
var (
    //go:embed whitefont.json
    WhiteFont_json []byte
    //go:embed whitefont.png
    WhiteFont_png []byte

    //go:embed blackfont.json
    BlackFont_json []byte
    //go:embed blackfont.png
    BlackFont_png []byte

    //go:embed small_asteroid.json 
    SmallAsteroid_json []byte
    //go:embed small_asteroid.png 
    SmallAsteroid_png []byte

    //go:embed medium_asteroid.json 
    MediumAsteroid_json []byte
    //go:embed medium_asteroid.png 
    MediumAsteroid_png []byte

    //go:embed large_asteroid.json 
    LargeAsteroid_json []byte
    //go:embed large_asteroid.png 
    LargeAsteroid_png []byte

    //go:embed small_explosion.json 
    SmallExplosion_json []byte
    //go:embed small_explosion.png 
    SmallExplosion_png []byte

    //go:embed simple_bullet.json 
    SimpleBullet_json []byte
    //go:embed simple_bullet.png
    SimpleBullet_png []byte
)
