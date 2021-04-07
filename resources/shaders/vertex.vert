#version 410
in vec2 vert;
in vec4 rotgroup;
in vec2 verttexcoord;

//Translation, window dimension scaling, rotation
uniform vec2 trans;
uniform mat2 dim;
uniform mat2 rot;
uniform vec2 rotcenter;

// Camera and zoom
uniform float zoom;
uniform vec2 cam;

out vec2 fragtexcoord;
void main(){
    // Set tex coords for frag shader
    fragtexcoord=verttexcoord;
    vec2 pos=vert;
    
    //Apply rotgroup rotation first, we want local changes then global changes to each vertex
    // vec2 rotcenter=vec2(rotgroup.x,rotgroup.y);
    // pos-=rotcenter;
    
    // mat2 rotmat=mat2(
        //     rotgroup.z,rotgroup.w,
        //     -rotgroup.w,rotgroup.z
    // );
    // pos=rotmat*pos;
    
    // pos+=rotcenter;
    
    // Apply uniform rotation
    pos=pos-rotcenter;
    pos=rot*pos;
    pos=pos+rotcenter;
    
    // Apply screen scaling from pixel coordinates
    vec2 I=vec2(1,1);
    pos=zoom*(dim*(pos+trans-cam))-I;
    
    gl_Position=vec4(pos,0.,1.);
}