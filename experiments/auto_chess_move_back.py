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

TURN_DURATION_MS = 1000
PAWN_KILL_REWARD = 4
STARTING_HEALTH = 10
BREACH_DAMAGE = 1
KING_CAPTURE_DAMAGE = 1
RAMP_END_ROUND = 120

PIECE_VALUE = {
    "P": 1,
    "N": 3,
    "B": 3,
    "R": 5,
    "Q": 9,
    "K": 100,
}
PIECE_LABEL = {
    "P": "pawn",
    "N": "knight",
    "B": "bishop",
    "R": "rook",
    "Q": "queen",
}

LIGHT_TILE = (238, 217, 182)
DARK_TILE = (181, 136, 99)
BG_COLOR = (22, 26, 32)
PANEL_COLOR = (35, 44, 58)
PANEL_BORDER = (16, 20, 27)
DEPLOY_HINT = (95, 160, 120)
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


@dataclass
class ForcedReturn:
    from_sq: tuple[int, int]
    to_sq: tuple[int, int]
    due_turn: int


class AutoChessMoveBack:
    def __init__(self) -> None:
        pygame.init()
        pygame.display.set_caption("Auto Chess Move Back")
        self.screen = pygame.display.set_mode((WINDOW_WIDTH, WINDOW_HEIGHT))
        self.clock = pygame.time.Clock()

        self.title_font = pygame.font.SysFont("georgia", 36, bold=True)
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
        y = 280
        for kind, label, cost in items:
            rect = pygame.Rect(PANEL_LEFT + 12, y, 250, 52)
            result.append(ShopItem(kind=kind, label=label, cost=cost, rect=rect))
            y += 66
        return result

    def reset(self) -> None:
        self.board: list[list[Optional[Piece]]] = [[None for _ in range(8)] for _ in range(8)]
        self.board[7][4] = Piece("white", "K")

        self.money = 30
        self.health = STARTING_HEALTH
        self.kills = 0
        self.turn_count = 0
        self.black_round = 0

        self.selected_shop_kind: Optional[str] = None
        self.pending_forced_return: Optional[ForcedReturn] = None

        self.game_over = False
        self.game_over_reason = ""
        self.current_turn = "black"
        self.status_text = "Auto move-back mode: white units retreat next turn."

        self.next_turn_ms = pygame.time.get_ticks() + TURN_DURATION_MS

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

    def find_white_king(self) -> Optional[tuple[int, int]]:
        for y in range(8):
            for x in range(8):
                piece = self.get_piece(x, y)
                if piece is not None and piece.color == "white" and piece.kind == "K":
                    return (x, y)
        return None

    def black_pawn_positions(self) -> list[tuple[int, int]]:
        out: list[tuple[int, int]] = []
        for y in range(8):
            for x in range(8):
                piece = self.get_piece(x, y)
                if piece is not None and piece.color == "black":
                    out.append((x, y))
        return out

    def square_from_mouse(self, mx: int, my: int) -> Optional[tuple[int, int]]:
        if not (BOARD_LEFT <= mx < BOARD_LEFT + BOARD_PIXELS and BOARD_TOP <= my < BOARD_TOP + BOARD_PIXELS):
            return None
        x = (mx - BOARD_LEFT) // TILE_SIZE
        y = (my - BOARD_TOP) // TILE_SIZE
        return int(x), int(y)

    def is_square_attacked_by_black_pawn(self, x: int, y: int) -> bool:
        for px in (x - 1, x + 1):
            py = y - 1
            if not self.in_bounds(px, py):
                continue
            piece = self.get_piece(px, py)
            if piece is not None and piece.color == "black" and piece.kind == "P":
                return True
        return False

    def legal_moves_for_white(self, x: int, y: int) -> set[tuple[int, int]]:
        piece = self.get_piece(x, y)
        if piece is None or piece.color != "white":
            return set()

        if piece.kind == "K":
            return self.king_moves(x, y, "white", avoid_black_pawn_attacks=True)
        if piece.kind == "N":
            return self.knight_moves(x, y, "white")
        if piece.kind == "R":
            return self.sliding_moves(x, y, "white", ((1, 0), (-1, 0), (0, 1), (0, -1)))
        if piece.kind == "B":
            return self.sliding_moves(x, y, "white", ((1, 1), (-1, 1), (1, -1), (-1, -1)))
        if piece.kind == "P":
            return self.white_pawn_moves(x, y)
        return set()

    def king_moves(
        self,
        x: int,
        y: int,
        color: str,
        avoid_black_pawn_attacks: bool = False,
    ) -> set[tuple[int, int]]:
        targets: set[tuple[int, int]] = set()
        for dy in (-1, 0, 1):
            for dx in (-1, 0, 1):
                if dx == 0 and dy == 0:
                    continue
                nx, ny = x + dx, y + dy
                if not self.in_bounds(nx, ny):
                    continue
                target = self.get_piece(nx, ny)
                if target is not None and target.color == color:
                    continue
                if avoid_black_pawn_attacks and self.is_square_attacked_by_black_pawn(nx, ny):
                    continue
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

    def black_pawn_moves(self, x: int, y: int) -> set[tuple[int, int]]:
        targets: set[tuple[int, int]] = set()
        ny = y + 1
        if self.in_bounds(x, ny) and self.get_piece(x, ny) is None:
            targets.add((x, ny))
        for nx in (x - 1, x + 1):
            if not self.in_bounds(nx, ny):
                continue
            target = self.get_piece(nx, ny)
            if target is not None and target.color == "white":
                targets.add((nx, ny))
        return targets

    def legal_moves_for_black(self, x: int, y: int) -> set[tuple[int, int]]:
        piece = self.get_piece(x, y)
        if piece is None or piece.color != "black":
            return set()

        if piece.kind == "P":
            return self.black_pawn_moves(x, y)
        if piece.kind == "N":
            return self.knight_moves(x, y, "black")
        if piece.kind == "B":
            return self.sliding_moves(x, y, "black", ((1, 1), (-1, 1), (1, -1), (-1, -1)))
        if piece.kind == "R":
            return self.sliding_moves(x, y, "black", ((1, 0), (-1, 0), (0, 1), (0, -1)))
        if piece.kind == "Q":
            return self.sliding_moves(
                x,
                y,
                "black",
                ((1, 1), (-1, 1), (1, -1), (-1, -1), (1, 0), (-1, 0), (0, 1), (0, -1)),
            )
        if piece.kind == "K":
            return self.king_moves(x, y, "black")
        return set()

    def black_pawn_attack_king_next_turn(self, pawn_x: int, pawn_y: int) -> bool:
        king = self.find_white_king()
        if king is None:
            return False
        kx, ky = king
        return ky == pawn_y + 1 and abs(kx - pawn_x) == 1

    def nearest_black_distance(self, x: int, y: int) -> Optional[int]:
        pawns = self.black_pawn_positions()
        if not pawns:
            return None
        return min(abs(x - px) + abs(y - py) for px, py in pawns)

    def score_white_move(
        self,
        piece: Piece,
        from_sq: tuple[int, int],
        to_sq: tuple[int, int],
        king_threatened: bool,
    ) -> float:
        fx, fy = from_sq
        tx, ty = to_sq
        score = self.rng.uniform(-0.08, 0.08)

        target = self.get_piece(tx, ty)
        if target is not None and target.color == "black":
            score += 120.0
            if self.black_pawn_attack_king_next_turn(tx, ty):
                score += 45.0

        before = self.nearest_black_distance(fx, fy)
        after = self.nearest_black_distance(tx, ty)
        if before is not None and after is not None:
            score += (before - after) * 5.0

        if piece.kind == "P":
            score += (fy - ty) * 6.0
        elif piece.kind in ("R", "B", "N"):
            score += max(0.0, 3.5 - abs(tx - 3.5))
        elif piece.kind == "K":
            if target is not None:
                score += 150.0
            score -= 12.0
            if king_threatened:
                score += 28.0
            if self.is_square_attacked_by_black_pawn(tx, ty):
                score -= 1000.0

        return score

    def choose_move_for_white_piece(
        self,
        piece: Piece,
        from_sq: tuple[int, int],
        legal_moves: list[tuple[int, int]],
        king_threatened: bool,
    ) -> Optional[tuple[int, int]]:
        if not legal_moves:
            return None

        best_score = float("-inf")
        best_moves: list[tuple[int, int]] = []
        for to_sq in legal_moves:
            score = self.score_white_move(piece, from_sq, to_sq, king_threatened)
            if score > best_score + 1e-9:
                best_score = score
                best_moves = [to_sq]
            elif abs(score - best_score) <= 1e-9:
                best_moves.append(to_sq)

        return self.rng.choice(best_moves) if best_moves else None

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

        self.set_piece(tx, ty, piece)
        self.set_piece(fx, fy, None)

        if piece.kind == "P" and ty == 0:
            self.set_piece(tx, ty, Piece("white", "R"))

    def white_turn(self) -> None:
        king = self.find_white_king()
        king_threatened = bool(king and self.is_square_attacked_by_black_pawn(king[0], king[1]))

        priority = {"P": 0, "N": 1, "B": 2, "R": 3, "K": 4}
        pieces: list[tuple[str, int, int]] = []
        for y in range(8):
            for x in range(8):
                piece = self.get_piece(x, y)
                if piece is not None and piece.color == "white":
                    pieces.append((piece.kind, x, y))
        pieces.sort(key=lambda entry: (priority.get(entry[0], 9), entry[2]))

        best_choice: Optional[tuple[float, tuple[int, int], tuple[int, int]]] = None
        for _, x, y in pieces:
            piece = self.get_piece(x, y)
            if piece is None or piece.color != "white":
                continue

            legal = list(self.legal_moves_for_white(x, y))
            chosen = self.choose_move_for_white_piece(piece, (x, y), legal, king_threatened)
            if chosen is None:
                continue

            score = self.score_white_move(piece, (x, y), chosen, king_threatened)
            if best_choice is None or score > best_choice[0] + 1e-9:
                best_choice = (score, (x, y), chosen)
            elif best_choice is not None and abs(score - best_choice[0]) <= 1e-9:
                if self.rng.random() < 0.5:
                    best_choice = (score, (x, y), chosen)

        if best_choice is None:
            self.status_text = "White turn: no useful moves."
        else:
            _, from_sq, to_sq = best_choice
            self.move_white_piece(from_sq, to_sq)
            self.pending_forced_return = ForcedReturn(
                from_sq=from_sq,
                to_sq=to_sq,
                due_turn=self.turn_count + 1,
            )
            self.status_text = "White turn: 1 auto move."

    def difficulty_progress(self) -> float:
        return min(1.0, self.black_round / RAMP_END_ROUND)

    def spawn_interval_for_progress(self, progress: float) -> int:
        if progress < 0.35:
            return 3
        if progress < 0.7:
            return 2
        return 1

    def spawn_weights_for_progress(self, progress: float) -> list[tuple[str, float]]:
        if progress < 0.2:
            return [("P", 1.0)]
        if progress < 0.4:
            return [("P", 0.78), ("N", 0.12), ("B", 0.10)]
        if progress < 0.6:
            return [("P", 0.50), ("N", 0.20), ("B", 0.20), ("R", 0.10)]
        if progress < 0.8:
            return [("P", 0.25), ("N", 0.18), ("B", 0.18), ("R", 0.19), ("Q", 0.20)]
        if progress < 0.95:
            return [("P", 0.08), ("N", 0.10), ("B", 0.10), ("R", 0.17), ("Q", 0.55)]
        return [("Q", 1.0)]

    def should_spawn_this_round(self) -> bool:
        progress = self.difficulty_progress()
        interval = self.spawn_interval_for_progress(progress)
        return self.black_round % interval == 0

    def choose_spawn_kind(self) -> str:
        weights = self.spawn_weights_for_progress(self.difficulty_progress())
        kinds = [kind for kind, _ in weights]
        probs = [weight for _, weight in weights]
        return self.rng.choices(kinds, weights=probs, k=1)[0]

    def spawn_black_piece(self, kind: str) -> bool:
        preferred_rows = [1, 0] if kind == "P" else [0, 1]
        for row in preferred_rows:
            candidates = [x for x in range(8) if self.get_piece(x, row) is None]
            if not candidates:
                continue
            x = self.rng.choice(candidates)
            self.set_piece(x, row, Piece("black", kind))
            return True

        fallback = [(x, y) for y in (0, 1) for x in range(8) if self.get_piece(x, y) is None]
        if not fallback:
            return False
        x, y = self.rng.choice(fallback)
        self.set_piece(x, y, Piece("black", kind))
        return True

    def try_spawn_black_piece_for_round(self) -> tuple[Optional[str], bool]:
        if not self.should_spawn_this_round():
            return None, False
        kind = self.choose_spawn_kind()
        if not self.spawn_black_piece(kind):
            return None, True
        return kind, True

    def score_black_move(
        self,
        piece: Piece,
        from_sq: tuple[int, int],
        to_sq: tuple[int, int],
    ) -> float:
        fx, fy = from_sq
        tx, ty = to_sq
        score = self.rng.uniform(-0.08, 0.08)

        target = self.get_piece(tx, ty)
        if target is not None and target.color == "white":
            score += PIECE_VALUE.get(target.kind, 1) * 45.0
            if target.kind == "K":
                score += 10000.0

        king = self.find_white_king()
        if king is not None:
            kx, ky = king
            before = abs(fx - kx) + abs(fy - ky)
            after = abs(tx - kx) + abs(ty - ky)
            score += (before - after) * 4.0

        if piece.kind == "P":
            score += (ty - fy) * 2.5
            if ty == 7:
                score += 8.0
        elif piece.kind in ("R", "Q"):
            score += max(0.0, ty - 3) * 0.8

        return score

    def choose_move_for_black_piece(
        self,
        piece: Piece,
        from_sq: tuple[int, int],
        legal_moves: list[tuple[int, int]],
    ) -> Optional[tuple[int, int]]:
        if not legal_moves:
            return None

        best_score = float("-inf")
        best_moves: list[tuple[int, int]] = []
        for to_sq in legal_moves:
            score = self.score_black_move(piece, from_sq, to_sq)
            if score > best_score + 1e-9:
                best_score = score
                best_moves = [to_sq]
            elif abs(score - best_score) <= 1e-9:
                best_moves.append(to_sq)
        return self.rng.choice(best_moves) if best_moves else None

    def move_black_pieces(self) -> tuple[int, int]:
        moved = 0
        king_captures = 0
        positions = [
            (x, y)
            for y in range(7, -1, -1)
            for x in range(8)
            if (self.get_piece(x, y) is not None and self.get_piece(x, y).color == "black")
        ]

        for x, y in positions:
            piece = self.get_piece(x, y)
            if piece is None or piece.color != "black":
                continue

            legal = list(self.legal_moves_for_black(x, y))
            chosen = self.choose_move_for_black_piece(piece, (x, y), legal)
            if chosen is None:
                continue

            tx, ty = chosen
            target = self.get_piece(tx, ty)
            if target is not None and target.color == "white" and target.kind == "K":
                # King capture deals damage but attacker is removed to avoid stacking breach damage.
                self.set_piece(tx, ty, None)
                self.set_piece(x, y, None)
                moved += 1
                king_captures += 1
                continue

            self.set_piece(tx, ty, piece)
            self.set_piece(x, y, None)
            moved += 1

        return moved, king_captures

    def black_turn(self) -> None:
        self.black_round += 1
        moved, king_captures = self.move_black_pieces()

        if king_captures > 0:
            damage = KING_CAPTURE_DAMAGE * king_captures
            label = (
                "Black captured your king."
                if king_captures == 1
                else f"Black captured your king {king_captures} times."
            )
            self.apply_damage(damage, label)
            if not self.game_over and self.find_white_king() is None and self.respawn_king():
                self.status_text = f"{self.status_text} King respawned."

        spawned_kind, spawn_attempted = self.try_spawn_black_piece_for_round()

        if self.game_over:
            return

        if spawn_attempted:
            if spawned_kind is None:
                spawn_text = "spawn blocked"
            else:
                spawn_text = f"spawned 1 {PIECE_LABEL[spawned_kind]}"
        else:
            spawn_text = "no spawn this round"

        if king_captures > 0:
            self.status_text = f"{self.status_text} {spawn_text}."
        else:
            self.status_text = f"Black round {self.black_round}: {moved} moves, {spawn_text}."

    def process_breaches(self) -> int:
        breached_files = [
            x
            for x in range(8)
            if (self.get_piece(x, 7) is not None and self.get_piece(x, 7).color == "black")
        ]
        for x in breached_files:
            self.set_piece(x, 7, None)

        breaches = len(breached_files)
        if breaches > 0:
            damage = BREACH_DAMAGE * breaches
            label = "A black pawn breached your back rank." if breaches == 1 else f"{breaches} black pawns breached your back rank."
            self.apply_damage(damage, label)
        return breaches

    def apply_damage(self, amount: int, reason: str) -> None:
        if amount <= 0 or self.game_over:
            return
        self.health = max(0, self.health - amount)
        if self.health <= 0:
            self.defeat(f"{reason} Health reached 0.")
            return
        self.status_text = f"{reason} -{amount} HP (health {self.health})."

    def respawn_king(self) -> bool:
        if self.find_white_king() is not None:
            return True
        preferred = [
            (4, 7), (3, 7), (5, 7), (2, 7), (6, 7), (1, 7), (7, 7), (0, 7),
            (4, 6), (3, 6), (5, 6), (2, 6), (6, 6), (1, 6), (7, 6), (0, 6),
        ]
        for x, y in preferred:
            if self.get_piece(x, y) is None:
                self.set_piece(x, y, Piece("white", "K"))
                return True
        return False

    def defeat(self, reason: str) -> None:
        self.game_over = True
        self.game_over_reason = reason
        self.selected_shop_kind = None

    def check_defeat_conditions(self) -> None:
        if self.health <= 0 and not self.game_over:
            self.defeat("Health reached 0.")

    def process_forced_return(self) -> None:
        if self.pending_forced_return is None:
            return
        if self.turn_count < self.pending_forced_return.due_turn:
            return

        ret = self.pending_forced_return
        self.pending_forced_return = None

        from_x, from_y = ret.from_sq
        to_x, to_y = ret.to_sq
        mover = self.get_piece(to_x, to_y)
        if mover is None or mover.color != "white":
            return

        target = self.get_piece(from_x, from_y)
        if target is not None and target.color == "black" and target.kind == "P":
            self.money += PAWN_KILL_REWARD
            self.kills += 1

        if target is not None and target.color == "white":
            # Retreat is mandatory: if ally blocks the base square, swap them.
            self.set_piece(to_x, to_y, target)
        else:
            self.set_piece(to_x, to_y, None)
        self.set_piece(from_x, from_y, mover)
        self.status_text = "Forced return executed."

    def resolve_turn(self) -> None:
        self.process_forced_return()

        if self.current_turn == "white":
            self.white_turn()
            self.current_turn = "black"
        else:
            self.black_turn()
            self.current_turn = "white"

        self.process_breaches()
        self.turn_count += 1
        self.check_defeat_conditions()

    def update(self) -> None:
        if self.game_over:
            return

        now = pygame.time.get_ticks()
        while now >= self.next_turn_ms and not self.game_over:
            self.resolve_turn()
            self.next_turn_ms += TURN_DURATION_MS

    def choose_shop_item(self, mouse_pos: tuple[int, int]) -> bool:
        mx, my = mouse_pos
        for item in self.shop_items:
            if not item.rect.collidepoint(mx, my):
                continue
            if self.money < item.cost:
                self.status_text = f"Need {item.cost} gold for {item.label}."
                return True
            self.selected_shop_kind = item.kind
            self.status_text = f"Placing {item.label}: click empty tile on rows 5-8."
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
            self.status_text = "Deploy only on rows 5-8."
            return
        if self.get_piece(x, y) is not None:
            self.status_text = "That tile is occupied."
            return
        if self.money < item.cost:
            self.status_text = "Not enough gold."
            self.selected_shop_kind = None
            return

        self.money -= item.cost
        self.set_piece(x, y, Piece("white", item.kind))
        self.selected_shop_kind = None
        self.status_text = f"Deployed {item.label}."

    def handle_click(self, mouse_pos: tuple[int, int]) -> None:
        if self.game_over:
            return

        if self.choose_shop_item(mouse_pos):
            return

        board_square = self.square_from_mouse(*mouse_pos)
        if board_square is None:
            return

        if self.selected_shop_kind is None:
            self.status_text = "Buy a piece from the shop first."
            return

        self.try_place_piece(board_square)

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
                pygame.draw.rect(self.screen, color, rect)

                if self.selected_shop_kind is not None and y >= 4 and self.get_piece(x, y) is None:
                    pygame.draw.rect(self.screen, DEPLOY_HINT, rect, width=2)

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

        health_color = DANGER_TEXT if self.health <= 3 else GOOD_TEXT
        health = self.text_font.render(f"Health: {self.health}", True, health_color)
        self.screen.blit(health, (PANEL_LEFT + 12, BOARD_TOP + 86))

        kills = self.small_font.render(f"Pawns killed: {self.kills}", True, MUTED_TEXT)
        self.screen.blit(kills, (PANEL_LEFT + 12, BOARD_TOP + 126))

        turn = self.small_font.render(f"Turn: {self.current_turn.title()}", True, TEXT_COLOR)
        self.screen.blit(turn, (PANEL_LEFT + 12, BOARD_TOP + 156))

        if self.game_over:
            timer_label = self.small_font.render("Next turn: stopped", True, DANGER_TEXT)
        else:
            seconds = max(0.0, (self.next_turn_ms - pygame.time.get_ticks()) / 1000)
            timer_label = self.small_font.render(f"Next turn in: {seconds:0.1f}s", True, TEXT_COLOR)
        self.screen.blit(timer_label, (PANEL_LEFT + 12, BOARD_TOP + 186))

        turns = self.small_font.render(f"Turns elapsed: {self.turn_count}", True, MUTED_TEXT)
        self.screen.blit(turns, (PANEL_LEFT + 12, BOARD_TOP + 216))

        spawn_interval = self.spawn_interval_for_progress(self.difficulty_progress())
        spawn_text = "every round" if spawn_interval == 1 else f"every {spawn_interval} rounds"
        round_info = self.small_font.render(
            f"Black rounds: {self.black_round}  |  Spawn: {spawn_text}",
            True,
            MUTED_TEXT,
        )
        self.screen.blit(round_info, (PANEL_LEFT + 12, BOARD_TOP + 244))

        for item in self.shop_items:
            affordable = self.money >= item.cost
            selected = self.selected_shop_kind == item.kind
            fill = (74, 126, 83) if selected else (56, 70, 90) if affordable else (66, 58, 58)
            pygame.draw.rect(self.screen, fill, item.rect, border_radius=8)
            pygame.draw.rect(self.screen, PANEL_BORDER, item.rect, width=2, border_radius=8)

            label = self.small_font.render(f"{item.label}  (${item.cost})", True, TEXT_COLOR)
            self.screen.blit(label, (item.rect.x + 12, item.rect.y + 15))

        instruction_lines = [
            "Auto-battle rules:",
            "- White and black alternate turns.",
            "- Each turn is exactly 1 second.",
            "- White auto-moves exactly one piece each turn.",
            "- The moved white piece must return next turn.",
            "- Forced return is free (does not consume a turn).",
            "- Black army moves automatically.",
            "- Start: pawn spawn every 3rd black round.",
            "- Endgame: one spawn every round, queens only.",
            f"- Start with {STARTING_HEALTH} health.",
            f"- Breach or king capture: -1 health.",
            f"- Kill pawn: +{PAWN_KILL_REWARD} gold.",
            "- You only buy and place pieces.",
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

        title = self.title_font.render("Auto Chess Move Back", True, TEXT_COLOR)
        self.screen.blit(title, (bar.x + 12, bar.y + 2))

        status_color = DANGER_TEXT if self.game_over else MUTED_TEXT
        status = self.small_font.render(self.status_text, True, status_color)
        self.screen.blit(status, (bar.x + 470, bar.y + 13))

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
    AutoChessMoveBack().run()
