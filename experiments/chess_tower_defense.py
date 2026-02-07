from __future__ import annotations

import random
from dataclasses import dataclass
from pathlib import Path
from typing import Optional

import pygame

BOARD_SIZE = 8
TILE_SIZE = 84
BOARD_PIXELS = BOARD_SIZE * TILE_SIZE

WINDOW_WIDTH = 1160
WINDOW_HEIGHT = 860

BOARD_LEFT = 42
BOARD_TOP = 86
PANEL_LEFT = BOARD_LEFT + BOARD_PIXELS + 36
FPS = 60

ENEMY_STEP_MS = 3000
PAWN_KILL_REWARD = 4

LIGHT_TILE = (238, 217, 182)
DARK_TILE = (181, 136, 99)
BG_COLOR = (22, 26, 32)
PANEL_COLOR = (35, 44, 58)
PANEL_BORDER = (16, 20, 27)
SELECT_COLOR = (90, 166, 120)
TARGET_COLOR = (91, 140, 196)
TEXT_COLOR = (236, 240, 245)
MUTED_TEXT = (179, 188, 202)
DANGER_TEXT = (244, 124, 124)
GOOD_TEXT = (135, 224, 162)

ATLAS_TILE = 16
PIECE_SIZE = TILE_SIZE - 16


@dataclass(frozen=True)
class Piece:
    color: str
    kind: str


@dataclass
class ShopItem:
    kind: str
    label: str
    cost: int
    rect: pygame.Rect


class ChessTowerDefense:
    def __init__(self) -> None:
        pygame.init()
        pygame.display.set_caption("Chess Tower Defense")
        self.screen = pygame.display.set_mode((WINDOW_WIDTH, WINDOW_HEIGHT))
        self.clock = pygame.time.Clock()

        self.title_font = pygame.font.SysFont("georgia", 38, bold=True)
        self.header_font = pygame.font.SysFont("georgia", 25, bold=True)
        self.text_font = pygame.font.SysFont("consolas", 24, bold=True)
        self.small_font = pygame.font.SysFont("consolas", 19)
        self.tiny_font = pygame.font.SysFont("consolas", 16)

        self.rng = random.Random()
        self.piece_sprites = self.load_piece_sprites()
        self.shop_items = self.create_shop()

        self.running = True
        self.reset()

    def create_shop(self) -> list[ShopItem]:
        items = [
            ("P", "Pawn", 5),
            ("N", "Knight", 12),
            ("B", "Bishop", 12),
            ("R", "Rook", 18),
        ]
        result: list[ShopItem] = []
        y = 268
        for kind, label, cost in items:
            rect = pygame.Rect(PANEL_LEFT + 12, y, 250, 52)
            result.append(ShopItem(kind=kind, label=label, cost=cost, rect=rect))
            y += 66
        return result

    def reset(self) -> None:
        self.board: list[list[Optional[Piece]]] = [[None for _ in range(8)] for _ in range(8)]
        self.board[7][4] = Piece("white", "K")

        self.money = 0
        self.kills = 0
        self.wave = 0

        self.selected_square: Optional[tuple[int, int]] = None
        self.legal_targets: set[tuple[int, int]] = set()
        self.selected_shop_kind: Optional[str] = None

        self.game_over = False
        self.game_over_reason = ""
        self.status_text = "Defend your king. Black pawns advance every 3 seconds."

        self.next_enemy_step_ms = pygame.time.get_ticks() + ENEMY_STEP_MS

    def load_piece_sprites(self) -> dict[tuple[str, str], pygame.Surface]:
        atlas_path = Path(__file__).resolve().parent.parent / "assets" / "roupiks" / "atlas.png"
        atlas = pygame.image.load(str(atlas_path)).convert_alpha()

        positions = {
            ("black", "K"): (1, 13),
            ("black", "Q"): (2, 13),
            ("black", "R"): (3, 13),
            ("black", "B"): (4, 13),
            ("black", "N"): (5, 13),
            ("black", "P"): (6, 13),
            ("white", "K"): (1, 14),
            ("white", "Q"): (2, 14),
            ("white", "R"): (3, 14),
            ("white", "B"): (4, 14),
            ("white", "N"): (5, 14),
            ("white", "P"): (6, 14),
        }

        sprites: dict[tuple[str, str], pygame.Surface] = {}
        for key, (tile_x, tile_y) in positions.items():
            src = pygame.Rect(tile_x * ATLAS_TILE, tile_y * ATLAS_TILE, ATLAS_TILE, ATLAS_TILE)
            tile = atlas.subsurface(src).copy()
            sprites[key] = pygame.transform.scale(tile, (PIECE_SIZE, PIECE_SIZE))
        return sprites

    def in_bounds(self, x: int, y: int) -> bool:
        return 0 <= x < 8 and 0 <= y < 8

    def get_piece(self, x: int, y: int) -> Optional[Piece]:
        return self.board[y][x]

    def set_piece(self, x: int, y: int, piece: Optional[Piece]) -> None:
        self.board[y][x] = piece

    def square_from_mouse(self, mx: int, my: int) -> Optional[tuple[int, int]]:
        if not (BOARD_LEFT <= mx < BOARD_LEFT + BOARD_PIXELS and BOARD_TOP <= my < BOARD_TOP + BOARD_PIXELS):
            return None
        x = (mx - BOARD_LEFT) // TILE_SIZE
        y = (my - BOARD_TOP) // TILE_SIZE
        return int(x), int(y)

    def legal_moves_for_white(self, x: int, y: int) -> set[tuple[int, int]]:
        piece = self.get_piece(x, y)
        if piece is None or piece.color != "white":
            return set()

        if piece.kind == "K":
            return self.king_moves(x, y, "white")
        if piece.kind == "N":
            return self.knight_moves(x, y, "white")
        if piece.kind == "R":
            return self.sliding_moves(x, y, "white", ((1, 0), (-1, 0), (0, 1), (0, -1)))
        if piece.kind == "B":
            return self.sliding_moves(x, y, "white", ((1, 1), (-1, 1), (1, -1), (-1, -1)))
        if piece.kind == "P":
            return self.white_pawn_moves(x, y)
        return set()

    def king_moves(self, x: int, y: int, color: str) -> set[tuple[int, int]]:
        targets: set[tuple[int, int]] = set()
        for dy in (-1, 0, 1):
            for dx in (-1, 0, 1):
                if dx == 0 and dy == 0:
                    continue
                nx, ny = x + dx, y + dy
                if not self.in_bounds(nx, ny):
                    continue
                target = self.get_piece(nx, ny)
                if target is None or target.color != color:
                    targets.add((nx, ny))
        return targets

    def knight_moves(self, x: int, y: int, color: str) -> set[tuple[int, int]]:
        targets: set[tuple[int, int]] = set()
        deltas = ((1, 2), (2, 1), (-1, 2), (-2, 1), (1, -2), (2, -1), (-1, -2), (-2, -1))
        for dx, dy in deltas:
            nx, ny = x + dx, y + dy
            if not self.in_bounds(nx, ny):
                continue
            target = self.get_piece(nx, ny)
            if target is None or target.color != color:
                targets.add((nx, ny))
        return targets

    def sliding_moves(
        self, x: int, y: int, color: str, directions: tuple[tuple[int, int], ...]
    ) -> set[tuple[int, int]]:
        targets: set[tuple[int, int]] = set()
        for dx, dy in directions:
            nx, ny = x + dx, y + dy
            while self.in_bounds(nx, ny):
                target = self.get_piece(nx, ny)
                if target is None:
                    targets.add((nx, ny))
                else:
                    if target.color != color:
                        targets.add((nx, ny))
                    break
                nx += dx
                ny += dy
        return targets

    def white_pawn_moves(self, x: int, y: int) -> set[tuple[int, int]]:
        targets: set[tuple[int, int]] = set()
        ny = y - 1
        if self.in_bounds(x, ny) and self.get_piece(x, ny) is None:
            targets.add((x, ny))
            if y == 6 and self.get_piece(x, y - 2) is None:
                targets.add((x, y - 2))
        for nx in (x - 1, x + 1):
            if not self.in_bounds(nx, ny):
                continue
            target = self.get_piece(nx, ny)
            if target is not None and target.color == "black":
                targets.add((nx, ny))
        return targets

    def move_white_piece(self, from_sq: tuple[int, int], to_sq: tuple[int, int]) -> None:
        fx, fy = from_sq
        tx, ty = to_sq
        piece = self.get_piece(fx, fy)
        if piece is None:
            return

        captured = self.get_piece(tx, ty)
        if captured is not None and captured.color == "black" and captured.kind == "P":
            self.money += PAWN_KILL_REWARD
            self.kills += 1
            self.status_text = f"+{PAWN_KILL_REWARD} gold for killing a pawn."

        self.set_piece(tx, ty, piece)
        self.set_piece(fx, fy, None)

        if piece.kind == "P" and ty == 0:
            self.set_piece(tx, ty, Piece("white", "R"))
            self.status_text = "Pawn promoted to rook."

    def choose_shop_item(self, mouse_pos: tuple[int, int]) -> bool:
        mx, my = mouse_pos
        for item in self.shop_items:
            if item.rect.collidepoint(mx, my):
                if self.money < item.cost:
                    self.status_text = f"Need {item.cost} gold for {item.label}."
                    return True
                self.selected_shop_kind = item.kind
                self.selected_square = None
                self.legal_targets.clear()
                self.status_text = f"Placing {item.label}: click an empty square on rows 5-8."
                return True
        return False

    def try_place_piece(self, square: tuple[int, int]) -> None:
        if self.selected_shop_kind is None:
            return

        x, y = square
        item = next((s for s in self.shop_items if s.kind == self.selected_shop_kind), None)
        if item is None:
            self.selected_shop_kind = None
            return

        if y < 4:
            self.status_text = "You can only deploy on your side (rows 5-8)."
            return
        if self.get_piece(x, y) is not None:
            self.status_text = "That square is occupied."
            return
        if self.money < item.cost:
            self.status_text = "Not enough gold."
            self.selected_shop_kind = None
            return

        self.money -= item.cost
        self.set_piece(x, y, Piece("white", item.kind))
        self.status_text = f"Deployed {item.label}."
        self.selected_shop_kind = None

    def handle_board_click(self, square: tuple[int, int]) -> None:
        if self.selected_shop_kind is not None:
            self.try_place_piece(square)
            return

        x, y = square
        clicked = self.get_piece(x, y)

        if self.selected_square is None:
            if clicked is not None and clicked.color == "white":
                self.selected_square = square
                self.legal_targets = self.legal_moves_for_white(x, y)
            return

        if square == self.selected_square:
            self.selected_square = None
            self.legal_targets.clear()
            return

        if square in self.legal_targets:
            self.move_white_piece(self.selected_square, square)
            self.selected_square = None
            self.legal_targets.clear()
            return

        if clicked is not None and clicked.color == "white":
            self.selected_square = square
            self.legal_targets = self.legal_moves_for_white(x, y)
        else:
            self.selected_square = None
            self.legal_targets.clear()

    def spawn_black_pawn(self) -> None:
        candidates = [x for x in range(8) if self.get_piece(x, 1) is None]
        if not candidates:
            return
        x = self.rng.choice(candidates)
        self.set_piece(x, 1, Piece("black", "P"))

    def move_black_pawns(self) -> None:
        for y in range(7, -1, -1):
            for x in range(8):
                piece = self.get_piece(x, y)
                if piece is None or piece.color != "black" or piece.kind != "P":
                    continue

                ny = y + 1
                if ny > 7:
                    continue

                captures: list[tuple[int, int, Piece]] = []
                for nx in (x - 1, x + 1):
                    if not self.in_bounds(nx, ny):
                        continue
                    target = self.get_piece(nx, ny)
                    if target is not None and target.color == "white":
                        captures.append((nx, ny, target))

                if captures:
                    captures.sort(key=lambda c: 0 if c[2].kind == "K" else 1)
                    tx, ty, taken = captures[0]
                    self.set_piece(tx, ty, piece)
                    self.set_piece(x, y, None)
                    if taken.kind == "K":
                        self.defeat("Your king was captured.")
                    continue

                if self.get_piece(x, ny) is None:
                    self.set_piece(x, ny, piece)
                    self.set_piece(x, y, None)

    def defeat(self, reason: str) -> None:
        self.game_over = True
        self.game_over_reason = reason
        self.selected_square = None
        self.legal_targets.clear()
        self.selected_shop_kind = None

    def check_defeat_conditions(self) -> None:
        king_alive = any(
            piece is not None and piece.color == "white" and piece.kind == "K"
            for row in self.board
            for piece in row
        )
        if not king_alive:
            self.defeat("Your king was captured.")
            return

        breached = any(
            self.get_piece(x, 7) is not None and self.get_piece(x, 7).color == "black"
            for x in range(8)
        )
        if breached:
            self.defeat("A black pawn breached your back rank.")

    def enemy_step(self) -> None:
        self.wave += 1
        self.spawn_black_pawn()
        self.move_black_pawns()
        self.check_defeat_conditions()
        if not self.game_over:
            self.status_text = f"Wave {self.wave}: black pawns advanced."

    def update(self) -> None:
        if self.game_over:
            return

        now = pygame.time.get_ticks()
        while now >= self.next_enemy_step_ms and not self.game_over:
            self.enemy_step()
            self.next_enemy_step_ms += ENEMY_STEP_MS

    def draw_piece(self, piece: Piece, center: tuple[int, int]) -> None:
        sprite = self.piece_sprites[(piece.color, piece.kind)]
        shadow = pygame.Surface((PIECE_SIZE, PIECE_SIZE), pygame.SRCALPHA)
        pygame.draw.ellipse(
            shadow,
            (0, 0, 0, 90),
            pygame.Rect((PIECE_SIZE - 14) // 2, PIECE_SIZE - 18, PIECE_SIZE - 14, 9),
        )
        shadow_rect = shadow.get_rect(center=(center[0], center[1] + PIECE_SIZE // 3))
        self.screen.blit(shadow, shadow_rect)
        self.screen.blit(sprite, sprite.get_rect(center=center))

    def draw_board(self) -> None:
        for y in range(8):
            for x in range(8):
                rect = pygame.Rect(
                    BOARD_LEFT + x * TILE_SIZE,
                    BOARD_TOP + y * TILE_SIZE,
                    TILE_SIZE,
                    TILE_SIZE,
                )
                color = DARK_TILE if (x + y) % 2 == 0 else LIGHT_TILE
                if self.selected_square == (x, y):
                    color = SELECT_COLOR
                elif (x, y) in self.legal_targets:
                    color = TARGET_COLOR
                pygame.draw.rect(self.screen, color, rect)

                if (x, y) in self.legal_targets and (x, y) != self.selected_square:
                    pygame.draw.circle(self.screen, (36, 54, 78), rect.center, TILE_SIZE // 9)

                piece = self.get_piece(x, y)
                if piece is not None:
                    self.draw_piece(piece, rect.center)

        pygame.draw.rect(
            self.screen,
            PANEL_BORDER,
            pygame.Rect(BOARD_LEFT, BOARD_TOP, BOARD_PIXELS, BOARD_PIXELS),
            width=4,
        )

        for idx, file_name in enumerate("abcdefgh"):
            text = self.tiny_font.render(file_name, True, MUTED_TEXT)
            tx = BOARD_LEFT + idx * TILE_SIZE + TILE_SIZE // 2 - text.get_width() // 2
            self.screen.blit(text, (tx, BOARD_TOP + BOARD_PIXELS + 6))

        for idx, rank_name in enumerate("87654321"):
            text = self.tiny_font.render(rank_name, True, MUTED_TEXT)
            ty = BOARD_TOP + idx * TILE_SIZE + TILE_SIZE // 2 - text.get_height() // 2
            self.screen.blit(text, (BOARD_LEFT - 16, ty))

    def draw_panel(self) -> None:
        panel = pygame.Rect(PANEL_LEFT, BOARD_TOP, WINDOW_WIDTH - PANEL_LEFT - 36, BOARD_PIXELS)
        pygame.draw.rect(self.screen, PANEL_COLOR, panel, border_radius=10)
        pygame.draw.rect(self.screen, PANEL_BORDER, panel, width=2, border_radius=10)

        title = self.header_font.render("Shop", True, TEXT_COLOR)
        self.screen.blit(title, (PANEL_LEFT + 12, BOARD_TOP + 10))

        money = self.text_font.render(f"Gold: {self.money}", True, GOOD_TEXT)
        self.screen.blit(money, (PANEL_LEFT + 12, BOARD_TOP + 48))

        kills = self.small_font.render(f"Pawns killed: {self.kills}", True, MUTED_TEXT)
        self.screen.blit(kills, (PANEL_LEFT + 12, BOARD_TOP + 86))

        if self.game_over:
            timer_label = self.small_font.render("Timer: stopped", True, DANGER_TEXT)
        else:
            seconds = max(0.0, (self.next_enemy_step_ms - pygame.time.get_ticks()) / 1000)
            timer_label = self.small_font.render(f"Next pawn step: {seconds:0.1f}s", True, TEXT_COLOR)
        self.screen.blit(timer_label, (PANEL_LEFT + 12, BOARD_TOP + 116))

        wave = self.small_font.render(f"Waves survived: {self.wave}", True, MUTED_TEXT)
        self.screen.blit(wave, (PANEL_LEFT + 12, BOARD_TOP + 144))

        for item in self.shop_items:
            affordable = self.money >= item.cost
            selected = self.selected_shop_kind == item.kind
            fill = (74, 126, 83) if selected else (56, 70, 90) if affordable else (66, 58, 58)
            pygame.draw.rect(self.screen, fill, item.rect, border_radius=8)
            pygame.draw.rect(self.screen, PANEL_BORDER, item.rect, width=2, border_radius=8)

            label = self.small_font.render(f"{item.label}  (${item.cost})", True, TEXT_COLOR)
            self.screen.blit(label, (item.rect.x + 12, item.rect.y + 15))

        instruction_lines = [
            "How to play:",
            "- Click white piece then destination.",
            "- Black pawns spawn/move every 3 seconds.",
            f"- Kill a pawn: +{PAWN_KILL_REWARD} gold.",
            "- Buy a piece, then click empty row 5-8 tile.",
            "- Survive as long as possible.",
            "R = restart, Esc = quit.",
        ]
        iy = BOARD_TOP + 548
        for line in instruction_lines:
            color = MUTED_TEXT if line.startswith("-") else TEXT_COLOR
            surf = self.tiny_font.render(line, True, color)
            self.screen.blit(surf, (PANEL_LEFT + 12, iy))
            iy += 21

    def draw_status_bar(self) -> None:
        bar = pygame.Rect(BOARD_LEFT, 26, WINDOW_WIDTH - 2 * BOARD_LEFT, 44)
        pygame.draw.rect(self.screen, PANEL_COLOR, bar, border_radius=8)
        pygame.draw.rect(self.screen, PANEL_BORDER, bar, width=2, border_radius=8)

        title = self.title_font.render("Chess Tower Defense", True, TEXT_COLOR)
        self.screen.blit(title, (bar.x + 12, bar.y + 2))

        status_color = DANGER_TEXT if self.game_over else MUTED_TEXT
        status = self.small_font.render(self.status_text, True, status_color)
        self.screen.blit(status, (bar.x + 460, bar.y + 13))

    def draw_game_over_overlay(self) -> None:
        if not self.game_over:
            return

        overlay = pygame.Surface((WINDOW_WIDTH, WINDOW_HEIGHT), pygame.SRCALPHA)
        overlay.fill((0, 0, 0, 142))
        self.screen.blit(overlay, (0, 0))

        msg = self.header_font.render("Game Over", True, DANGER_TEXT)
        msg_rect = msg.get_rect(center=(BOARD_LEFT + BOARD_PIXELS // 2, BOARD_TOP + BOARD_PIXELS // 2 - 18))
        self.screen.blit(msg, msg_rect)

        reason = self.small_font.render(self.game_over_reason, True, TEXT_COLOR)
        reason_rect = reason.get_rect(center=(BOARD_LEFT + BOARD_PIXELS // 2, BOARD_TOP + BOARD_PIXELS // 2 + 14))
        self.screen.blit(reason, reason_rect)

        hint = self.small_font.render("Press R to restart", True, MUTED_TEXT)
        hint_rect = hint.get_rect(center=(BOARD_LEFT + BOARD_PIXELS // 2, BOARD_TOP + BOARD_PIXELS // 2 + 44))
        self.screen.blit(hint, hint_rect)

    def draw(self) -> None:
        self.screen.fill(BG_COLOR)
        self.draw_status_bar()
        self.draw_board()
        self.draw_panel()
        self.draw_game_over_overlay()
        pygame.display.flip()

    def handle_click(self, mouse_pos: tuple[int, int]) -> None:
        if self.game_over:
            return

        if self.choose_shop_item(mouse_pos):
            return

        square = self.square_from_mouse(*mouse_pos)
        if square is None:
            return
        self.handle_board_click(square)

    def run(self) -> None:
        while self.running:
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    self.running = False
                elif event.type == pygame.KEYDOWN:
                    if event.key == pygame.K_ESCAPE:
                        self.running = False
                    elif event.key == pygame.K_r:
                        self.reset()
                elif event.type == pygame.MOUSEBUTTONDOWN and event.button == 1:
                    self.handle_click(event.pos)

            self.update()
            self.draw()
            self.clock.tick(FPS)

        pygame.quit()


if __name__ == "__main__":
    ChessTowerDefense().run()
