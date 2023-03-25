# Orange Character (Others to come)
---
I am loading all sprites associated to the characters movement in a 3x3 matrix. 
Each direction has 3 sprites:  

0. Idle state of direction
1. Animation state A for direction  
2. Animation state A for direction  

|   ↖   |   ⬆   |   ↗   |
| :---: | :---: | :---: |
|   ⬅   |   ?   |   ➡   |
|   ↙   |   ⬇   |   ↘   |

## Update Function
***Update gets called every tick (about $60TPS$)***  
   
Counter increments each time and flips flag **move** on and off for animation in the draw function.  

WASD keys are mapped to vectors with an x and y value that ranges from -1 to 1.
I'm adding 1 on lines 52/53 to map to the directions matrix [y][x].   
The x/y positon of the character is being multiplied by the characters speed.  
  
  
## Draw Function
***Draw gets called every frame (for me about $120FPS$)***  

This function only draws the sprite on the screen. 
