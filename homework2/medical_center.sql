CREATE TABLE IF NOT EXISTS "clinics" (
	"id" serial NOT NULL UNIQUE,
	"city" varchar(255) NOT NULL,
	"address" varchar(255) NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "doctors" (
	"id" serial NOT NULL UNIQUE,
	"name" varchar(255) NOT NULL,
	"surname" varchar(255) NOT NULL,
	"patronymic" varchar(255),
	"sex" varchar(255) NOT NULL,
	"speciality" varchar(255) NOT NULL,
	"is_active" boolean NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "patients" (
	"id" serial NOT NULL UNIQUE,
	"name" varchar(255) NOT NULL,
	"surname" varchar(255) NOT NULL,
	"patronymic" varchar(255),
	"sex" varchar(255) NOT NULL,
	"birthdate" date NOT NULL,
	"phone" varchar(255) NOT NULL,
	"email" bigint,
	"registrated_date" timestamp with time zone NOT NULL DEFAULT 'now()',
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "bookings" (
	"id" serial NOT NULL UNIQUE,
	"clinic_id" bigint NOT NULL,
	"patinent_id" bigint NOT NULL,
	"doctor_id" bigint NOT NULL,
	"appointment_time" timestamp with time zone NOT NULL,
	"booking_time" timestamp with time zone NOT NULL,
	"is_planned" boolean NOT NULL DEFAULT true,
	"is_approved" boolean NOT NULL DEFAULT false,
	"is_finished" boolean NOT NULL DEFAULT false,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "medical_history" (
	"id" serial NOT NULL UNIQUE,
	"clinic_id" bigint NOT NULL,
	"patient_id" bigint NOT NULL,
	"doctor_id" bigint NOT NULL,
	"booking_id" bigint NOT NULL,
	"appointment_date" date NOT NULL,
	"diagnostics" varchar(255) NOT NULL,
	"recommendations" varchar(255) NOT NULL,
	"record" varchar(255) NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "calls" (
	"id" serial NOT NULL UNIQUE,
	"phone" bigint NOT NULL,
	"call_time" timestamp with time zone NOT NULL,
	"booking_id" bigint NOT NULL,
	PRIMARY KEY ("id")
);




ALTER TABLE "bookings" ADD CONSTRAINT "bookings_fk1" FOREIGN KEY ("clinic_id") REFERENCES "clinics"("id");

ALTER TABLE "bookings" ADD CONSTRAINT "bookings_fk2" FOREIGN KEY ("patinent_id") REFERENCES "patients"("id");

ALTER TABLE "bookings" ADD CONSTRAINT "bookings_fk3" FOREIGN KEY ("doctor_id") REFERENCES "doctors"("id");
ALTER TABLE "medical_history" ADD CONSTRAINT "medical_history_fk1" FOREIGN KEY ("clinic_id") REFERENCES "clinics"("id");

ALTER TABLE "medical_history" ADD CONSTRAINT "medical_history_fk2" FOREIGN KEY ("patient_id") REFERENCES "patients"("id");

ALTER TABLE "medical_history" ADD CONSTRAINT "medical_history_fk3" FOREIGN KEY ("doctor_id") REFERENCES "doctors"("id");

ALTER TABLE "medical_history" ADD CONSTRAINT "medical_history_fk4" FOREIGN KEY ("booking_id") REFERENCES "bookings"("id");
ALTER TABLE "calls" ADD CONSTRAINT "calls_fk3" FOREIGN KEY ("booking_id") REFERENCES "bookings"("id");