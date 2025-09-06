# Dungeons & Quests Implementation Todo List

## Stage 1: Domain Layer Rename

- [ ] Commit 1: `refactor(domain): Task → Quest entities and interfaces`
  - [ ] Create `internal/domain/entity/quest.go`
  - [ ] Create `internal/domain/entity/dungeon.go`
  - [ ] Create `internal/domain/entity/quest_completion.go`
  - [ ] Update `internal/domain/entity/user.go`
  - [ ] Update `internal/domain/entity/shop_item.go`
  - [ ] Delete `internal/domain/entity/task.go`

## Stage 2: Repository Layer Rename

- [ ] Commit 2: `refactor(ports): TaskRepository → QuestRepository`
  - [ ] Update `internal/ports/repositories.go`

- [ ] Commit 3: `refactor(infra): rename repository implementations`
  - [ ] Rename `internal/infra/postgres/task_repository.go` → `quest_repository.go`
  - [ ] Update repository implementations

## Stage 3: Service Layer Rename

- [ ] Commit 4: `refactor(usecase): TaskService → QuestService`
  - [ ] Rename `internal/usecase/task_service.go` → `quest_service.go`
  - [ ] Update service methods and signatures
  - [ ] Rename test files

## Stage 4: HTTP Layer Rename

- [ ] Commit 5: `refactor(infra/http): /tasks → /quests route names`
  - [ ] Rename `internal/infra/http/task_handler.go` → `quest_handler.go`
  - [ ] Update `internal/infra/http/server.go`
  - [ ] Update DTOs

## Stage 5: Add New Domain Services

- [ ] Commit 6: `feat(domain): add Dungeon, DungeonMember, QuestCompletion (types only)`
  - [ ] Create `internal/usecase/dungeon_service.go`
  - [ ] Create repository implementations (stubs)

## Stage 6: Update Tests & Fixtures

- [ ] Commit 7: `refactor(infra/tg): replace chat_id references with Dungeon plumbing`
  - [ ] Update Telegram integration touchpoints

- [ ] Commit 8: `refactor(test): update fixtures and tests for quests`
  - [ ] Rename `test/fixtures/builders/task_builder.go` → `quest_builder.go`
  - [ ] Update test files

## Stage 7: Documentation Updates

- [ ] Commit 9: `docs: update names in README/tech overview`
  - [ ] Update `README.md`
  - [ ] Update `TECHNICAL_OVERVIEW.md`
  - [ ] Rename and update `docs/api/services/task_service.md` → `quest_service.md`

## Stage 8: Database Migration

- [ ] Commit 10: `feat(db): add migration for dungeons and quests`
  - [ ] Create `internal/infra/postgres/migrations/007_dungeons_and_quests.sql`

## Final Verification

- [ ] All `Task` types/interfaces renamed to `Quest`
- [ ] New domain types exist: `Dungeon`, `DungeonMember`, `QuestCompletion`
- [ ] `Quest` contains all MVP scoring fields
- [ ] Repository interfaces updated with new methods
- [ ] Service layer uses new names and signatures
- [ ] API routes use `/quests` and `/dungeons`
- [ ] JSON responses use new field names
- [ ] Tests updated with new builders
- [ ] Documentation reflects new terminology
- [ ] Database migration ready (fresh start)
- [ ] Build passes with all tests green