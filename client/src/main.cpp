#include "raylib.h"

int main() {
    InitWindow(800, 600, "Coin Collector");
    SetTargetFPS(60);

    Vector2 player = {400, 300};

    while (!WindowShouldClose()) {
        float dt = GetFrameTime();

        if (IsKeyDown(KEY_RIGHT)) player.x += 200 * dt;
        if (IsKeyDown(KEY_LEFT)) player.x -= 200 * dt;
        if (IsKeyDown(KEY_UP)) player.y -= 200 * dt;
        if (IsKeyDown(KEY_DOWN)) player.y += 200 * dt;

        BeginDrawing();
        ClearBackground(BLACK);
        DrawCircleV(player, 20, GREEN);
        DrawText("Raylib test running", 10, 10, 20, WHITE);
        EndDrawing();
    }

    CloseWindow();
    return 0;
}
