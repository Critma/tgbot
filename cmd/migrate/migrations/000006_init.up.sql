ALTER TABLE "reminders"
ADD COLUMN "task_id" text not null,
ADD COLUMN "task_queue" text not null;