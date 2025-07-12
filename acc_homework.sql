--
-- PostgreSQL database dump
--

-- Dumped from database version 14.15 (Homebrew)
-- Dumped by pg_dump version 14.15 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: assignments; Type: TABLE; Schema: public; Owner: williamfotso
--

CREATE TABLE public.assignments (
    id integer NOT NULL,
    course_code character varying,
    type character varying,
    deadline timestamp without time zone,
    title character varying,
    todo character varying,
    notion_id character varying,
    link text,
    status text DEFAULT 'default'::text
);


ALTER TABLE public.assignments OWNER TO williamfotso;

--
-- Name: assignements_id_seq; Type: SEQUENCE; Schema: public; Owner: williamfotso
--

ALTER TABLE public.assignments ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
    SEQUENCE NAME public.assignements_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: courses; Type: TABLE; Schema: public; Owner: williamfotso
--

CREATE TABLE public.courses (
    name character varying(40) NOT NULL,
    code character varying NOT NULL,
    notion_id text NOT NULL,
    duration text,
    room_number text,
    id integer NOT NULL
);


ALTER TABLE public.courses OWNER TO williamfotso;

--
-- Name: courses_id_seq; Type: SEQUENCE; Schema: public; Owner: williamfotso
--

CREATE SEQUENCE public.courses_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.courses_id_seq OWNER TO williamfotso;

--
-- Name: courses_id_seq1; Type: SEQUENCE; Schema: public; Owner: williamfotso
--

CREATE SEQUENCE public.courses_id_seq1
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.courses_id_seq1 OWNER TO williamfotso;

--
-- Name: courses_id_seq1; Type: SEQUENCE OWNED BY; Schema: public; Owner: williamfotso
--

ALTER SEQUENCE public.courses_id_seq1 OWNED BY public.courses.id;


--
-- Name: status; Type: TABLE; Schema: public; Owner: williamfotso
--

CREATE TABLE public.status (
    id text,
    name text,
    color text
);


ALTER TABLE public.status OWNER TO williamfotso;

--
-- Name: type; Type: TABLE; Schema: public; Owner: williamfotso
--

CREATE TABLE public.type (
    id text,
    color text NOT NULL,
    name text NOT NULL
);


ALTER TABLE public.type OWNER TO williamfotso;

--
-- Name: courses id; Type: DEFAULT; Schema: public; Owner: williamfotso
--

ALTER TABLE ONLY public.courses ALTER COLUMN id SET DEFAULT nextval('public.courses_id_seq1'::regclass);


--
-- Data for Name: assignments; Type: TABLE DATA; Schema: public; Owner: williamfotso
--

COPY public.assignments (id, course_code, type, deadline, title, todo, notion_id, link, status) FROM stdin;
75	GOVT-2305	HW	2025-07-08 00:00:00	Self Intro	use this forum to introduce ourselves to one another	22940a21-a7e3-8169-8f22-f1d3e8e0f709	https://acconline.austincc.edu/ultra/stream	done
77	GOVT-2305	HW	2025-07-10 00:00:00	Quiz 2	will be over Chapters 3&4	22940a21-a7e3-8114-879e-c49d7e1ee7e3	https://acconline.austincc.edu/ultra/courses/_942344_1/outline/assessment/_29676237_1/overview?courseId=_942344_1	default
78	ENGL-2311	HW	2025-07-09 00:00:00	Read Chapter 14: Completing Business Proposals and Reports/w questions (Cardon)	Do Read Chapter 14: Completing Business Proposals and Reports/w questions (Cardon)	22a40a21-a7e3-8175-b404-def89dc74553	https://acconline.austincc.edu/ultra/courses/_941962_1/outline/lti/launch?courseId=_941962_1&contentId=_29449771_1	
28	COSC-1301	Exam	2025-06-20 00:00:00	Concepts Exam Chapters 1-3	Concepts Exam Chapters 1-3	20d40a21-a7e3-81a4-90a6-f97187add2b7	https://acconline.austincc.edu/ultra/courses/_941496_1/outline/lti/launch?courseId=_941496_1&contentId=_28852683_1	done
70	ENGL-2311	HW	2025-07-09 00:00:00	Read Chapter 13: Research and Planning for Business Proposals and Reports/w questions (Cardon)	Do Read Chapter 13: Research and Planning for Business Proposals and Reports/w questions (Cardon)	22640a21-a7e3-817a-9315-d721db9fca45	https://acconline.austincc.edu/ultra/courses/_941962_1/grades/lti/launch?courseId=_941962_1&contentId=_29449784_1	default
73	ENGL-2311	Exam	2025-07-10 00:00:00	Powerpoint assignment	Powerpoint assignment	22840a21-a7e3-8064-8e81-d0e51087e5b1	https://acconline.austincc.edu/ultra/courses/_941962_1/announcements/announcement-detail?courseId=_941962_1&announcementId=_2610853_1	default
57	COSC-1301	Exam	2025-07-04 00:00:00	PowerPoint Modules 1-3 - SAM Capstone Project 1a	Do PowerPoint Modules 1-3 - SAM Capstone Project 1a	21f40a21-a7e3-8142-b45e-d3ff2c28c3a1	https://acconline.austincc.edu/ultra/courses/_941496_1/outline/lti/launch?courseId=_941496_1&contentId=_28852612_1	done
62	COSC-1301	HW	2025-07-18 00:00:00	Access Module 02 - SAM Project 1	Do access Module 02 - SAM Project 1	21f40a21-a7e3-8120-9897-e99ad7711b4d	https://acconline.austincc.edu/ultra/courses/_941496_1/outline/lti/launch?courseId=_941496_1&contentId=_28852892_1	default
66	ENGL-2311	HW	2025-07-02 00:00:00	Public Speaking Skills: Language and Delivery	Do public Speaking Skills: Language and Delivery	22340a21-a7e3-81e1-b28c-f8a9bdcc8b28	ine.austincc.edu/ultra/courses/_941962_1/outline/lti/launch?courseId=_941962_1&contentId=_29449824_1	done
61	COSC-1301	HW	2025-07-18 00:00:00	Access Module 01 - SAM Project 1	Do access Module 01 - SAM Project 1	21f40a21-a7e3-8166-ab0f-d37bbaab8cfb	https://acconline.austincc.edu/ultra/courses/_941496_1/outline/lti/launch?courseId=_941496_1&contentId=_28852868_1	default
65	ENGL-2311	HW	2025-07-02 00:00:00	Read: Delivering Presentations/w questions (Cardon)	DO Read: Delivering Presentations/w questions (Cardon)	22340a21-a7e3-8179-a4dc-f85a15318d8e	https://acconline.austincc.edu/ultra/courses/_941962_1/outline/lti/launch?courseId=_941962_1&contentId=_29449821_1	done
55	ENGL-2311	HW	2025-06-30 00:00:00	Public Speaking Skills: Speaking with Confidence	Do Public Speaking Skills: Speaking with Confidence	21f40a21-a7e3-8125-9aeb-e6580d478877	https://acconline.austincc.edu/ultra/courses/_941962_1/outline/lti/launch?courseId=_941962_1&contentId=_29449825_1	done
36	HIST-1301	HW	2025-06-21 00:00:00	KEY TERM QUIZ #1 {CH 1 - 4}	KEY TERM QUIZ #1 {CH 1 - 4}	21040a21-a7e3-81e2-ae3c-dd7492d12c0b	https://acconline.austincc.edu/ultra/courses/_942453_1/outline/edit/document/_28916920_1?courseId=_942453_1&view=content&state=view	done
27	COSC-1301	Exam	2025-06-15 00:00:00	PowerPoint Module 01 - SAM Project 1	PowerPoint Module 01 - SAM Project 1	20d40a21-a7e3-8132-ad58-f97690e18059	https://acconline.austincc.edu/ultra/courses/_941496_1/grades/lti/launch?courseId=_941496_1&contentId=_28852786_1	done
53	ENGL-2311	HW	2025-07-01 00:00:00	Quiz: Chapters 4,5,6 (Finkelstein)	Do Quiz: Chapters 4,5,6 (Finkelstein)	21f40a21-a7e3-8133-9a86-d1ab12b0e231	https://acconline.austincc.edu/ultra/courses/_941962_1/outline/lti/launch?courseId=_941962_1&contentId=_29685410_1	done
41	HIST-1301	HW	2025-07-19 00:00:00	KEY TERM QUIZ #3 {CH 9 - 12}	KEY TERM QUIZ #3 {CH 9 - 12}	21040a21-a7e3-8154-8b73-d6919dcdc4b5		default
42	HIST-1301	Exam	2025-07-19 00:00:00	EXAM #3 {CH 9 - 12}	EXAM #3 {CH 9 - 12}	21040a21-a7e3-811e-a9bb-d5e00643b78c		default
44	HIST-1301	HW	2025-08-02 00:00:00	KEY TERM QUIZ #4 {CH 13 - 16}	KEY TERM QUIZ #4 {CH 13 - 16}	21040a21-a7e3-8100-8f99-c9fd098dd644		default
45	HIST-1301	Exam	2025-08-02 00:00:00	EXAM #4 {CH 13 - 16}	EXAM #4 {CH 13 - 16}	21040a21-a7e3-81b1-a98f-fb8aec04ef0f		default
40	HIST-1301	HW	2025-07-12 00:00:00	Book Review Outline(s) 1 & 2 & Bibliography Page(s)	Book Review Outline(s) 1 & 2 & Bibliography Page(s)	21040a21-a7e3-81cf-b7ed-ed3cb27b292c	https://acconline.austincc.edu/ultra/courses/_942453_1/outline/edit/document/_28916905_1?courseId=_942453_1&view=content&state=view	default
43	HIST-1301	HW	2025-07-26 00:00:00	Book Review Final Drafts	Book Review Final Drafts	21040a21-a7e3-8190-96c8-c30a3e22e14e	https://acconline.austincc.edu/ultra/courses/_942453_1/outline/edit/document/_28916901_1?courseId=_942453_1&view=content&state=view	default
67	HIST-1301	HW	2025-07-03 00:00:00	Exam #1 extra credit	write a paragraph about the following topics that were missed on the exam. For each paragraph , I will give you 2 points back on your exam. Please note that a paragraph is 5 -7 sentences - minimum.\n\nPlease respond directly to this email.\n\n	22340a21-a7e3-81df-8c27-d17e90f56072	https://acconline.austincc.edu/ultra/stream	default
63	COSC-1301	HW	2025-07-25 00:00:00	Access Module 03 - SAM Project 1	Do access Module 03 - SAM Project 1	22040a21-a7e3-81da-84a4-ecebb71e4ad5	https://acconline.austincc.edu/ultra/courses/_941496_1/outline/lti/launch?courseId=_941496_1&contentId=_28852843_1	default
68	ENGL-2311	HW	2025-07-16 00:00:00	Instructions assignment	Write the instructions based on the images in the instruction.pdf file	22340a21-a7e3-81a0-900a-c2251e19bc9e	https://acconline.austincc.edu/ultra/courses/_941962_1/announcements/announcement-detail?courseId=_941962_1&announcementId=_2612050_1	default
58	COSC-1301	HW	2025-07-04 00:00:00	Excel Module 02 - SAM Project 1	Do excel Module 02 - SAM Project 1	21f40a21-a7e3-81af-9d55-fc47f22670cd	https://acconline.austincc.edu/ultra/courses/_941496_1/outline/lti/launch?courseId=_941496_1&contentId=_28852810_1	done
25	HIST-1301	Exam	2025-06-14 00:00:00	History orientation	History orientation	20d40a21-a7e3-8119-ac61-d27ad3cc0b01	https://acconline.austincc.edu/ultra/courses/_942453_1/outline/edit/document/_28916899_1?courseId=_942453_1&view=content&state=view	done
33	GOVT-2305	HW	2025-07-07 00:00:00	Start of Governement class	Do 1st Week Homework before 1st class	20f40a21-a7e3-81d9-8792-eea8aca6297a	https://acconline.austincc.edu/ultra/course	done
39	HIST-1301	Exam	2025-07-05 00:00:00	EXAM #2 {CH 5 - 8}	EXAM #2 {CH 5 - 8}	21040a21-a7e3-81af-8155-d67224104880	https://acconline.austincc.edu/ultra/courses/_942453_1/outline/assessment/_28916912_1/overview?courseId=_942453_1	done
38	HIST-1301	Exam	2025-07-05 00:00:00	KEY TERM QUIZ #2 {CH 5 - 8}	KEY TERM QUIZ #2 {CH 5 - 8}	21040a21-a7e3-8178-89a2-d7a6430b54f6	https://acconline.austincc.edu/ultra/courses/_942453_1/outline/assessment/_28916922_1/overview?courseId=_942453_1	done
52	ENGL-2311	Exam	2025-07-08 00:00:00	Powerpoint draft and rehearsal	- Finish draft of powerpoint presentation\n- Rehearse presentation	21c40a21-a7e3-81de-a3e7-f9c12f000aae	https://acconline.austincc.edu/ultra/courses/_941962_1/announcements/announcement-detail?courseId=_941962_1&announcementId=_2610853_1	start
64	COSC-1301	Exam	2025-08-01 00:00:00	Access Modules 1-3 - SAM Capstone Project 1a	Do access Modules 1-3 - SAM Capstone Project 1a	22040a21-a7e3-8185-84c5-f626c0411b9e	https://acconline.austincc.edu/ultra/courses/_941496_1/outline/lti/launch?courseId=_941496_1&contentId=_28852614_1	default
34	HIST-1301	HW	2025-06-14 00:00:00	BOOK CHOICES 1 & 2	BOOK CHOICES 1 & 2	21040a21-a7e3-812c-8f7f-c1139e27d138		done
35	HIST-1301	Exam	2025-06-17 00:00:00	MAP EXAM	MAP EXAM	21040a21-a7e3-8117-a7e1-d0c6ea365164		done
1	MATH-2412	HW	2025-03-28 00:00:00	ALEKS HW	Do missing aleks	1c240a21-a7e3-8165-af5d-df76cd31779d		done
2	ITNW-1325	Exam	2025-03-28 00:00:00	Network Exam	Revise the exam for next friday	1c240a21-a7e3-8131-abc0-f572b9ac99d5		done
3	MATH-2412	HW	2025-04-01 00:00:00	HomeWork 5	Do HW 5 (subject on BlackBoard)	1c540a21-a7e3-817e-8266-c33c70bd8773		done
4	ITSY-1300	HW	2025-04-11 00:00:00	Subject Term Paper	Find and Submit subject for security term paper	1c740a21-a7e3-81ee-acfb-c28e81d660e8		done
30	ENGL-2311	HW	2025-06-10 00:00:00	Role Play Assignment	Role Play Assignment	20d40a21-a7e3-81b0-8db1-c874616b30f6	https://acconline.austincc.edu/ultra/courses/_941962_1/outline/lti/launch?courseId=_941962_1&contentId=_29449763_1	done
7	COSC-2436	HW	2025-04-11 00:00:00	Prog Quiz	Do Programming Quiz 8 to Quiz Interlude 5	1ce40a21-a7e3-815e-8a75-fb3dbf152283		done
6	COSC-2436	HW	2025-04-11 00:00:00	Lab #3 Linked Lists	Start Lab #3	1c740a21-a7e3-8169-991e-e17682be5cfd		done
9	ITSY-1300	Exam	2025-04-30 00:00:00	Term Paper (Final)	Write and Submit Final Term paper on HomeLabs	1d140a21-a7e3-815f-aa5e-f29beb3b875e		done
10	MATH-2412	HW	2025-04-13 00:00:00	Quiz 7	Do Quiz 7 ( Vectors, Dot Product )	1d140a21-a7e3-8171-870e-e5294d335d86		done
8	ITSY-1300	Exam	2025-04-18 00:00:00	Term Paper (Draft)	Write a draft of each sections (1-6)	1d140a21-a7e3-8129-bfa9-e253a6712ec4		done
24	ENGL-2311	HW	2025-06-08 00:00:00	Quiz for Sunday	Quiz 1 & 2, Finish 1st Module	20a40a21-a7e3-81f3-ba03-f9b3612d3bea		done
46	COSC-1301	Exam	2025-06-20 00:00:00	Word Modules 1-3 - SAM Capstone Project 1a	Word Modules 1-3 - SAM Capstone Project 1a	21340a21-a7e3-8190-8fff-edbbde6c01fb	https://acconline.austincc.edu/ultra/courses/_941496_1/grades/lti/launch\\?courseId\\=_941496_1\\&contentId\\=_28852834_1	done
37	HIST-1301	Exam	2025-06-21 00:00:00	EXAM #1 {CH 1 – 4}	EXAM #1 {CH 1 – 4}	21040a21-a7e3-81de-bd12-c44627300540	https://acconline.austincc.edu/ultra/courses/_942453_1/outline/edit/document/_28916910_1?courseId=_942453_1&view=content&state=view	done
49	COSC-1301	HW	2025-06-27 00:00:00	Excel Module 01 - -SAM Project 1	Excel Module 01 - -SAM Project 1	21a40a21-a7e3-8139-bc08-f270e995514a	https://acconline.austincc.edu/ultra/courses/_941496_1/outline/lti/launch?courseId=_941496_1&contentId=_28852876_1	done
48	COSC-1301	Exam	2025-06-27 00:00:00	PowerPoint Module 03 - SAM Project 1	PowerPoint Module 03 - SAM Project 1	21a40a21-a7e3-81bf-9aa1-fd8ad85f5f5f	https://acconline.austincc.edu/ultra/courses/_941496_1/outline/lti/launch\\?courseId\\=_941496_1\\&contentId\\=_28852802_1	done
47	COSC-1301	Exam	2025-06-20 00:00:00	PowerPoint Module 02 - SAM Project 1	PowerPoint Module 02 - SAM Project 1	21340a21-a7e3-81fb-a03c-c363bb0dcdb4	https://acconline.austincc.edu/ultra/courses/_941496_1/grades/lti/launch\\?courseId\\=_941496_1\\&contentId\\=_28852884_1	done
26	COSC-1301	Exam	2025-06-15 00:00:00	Word Module 03 - SAM Project 1	Word Module 03 - SAM Project 1	20d40a21-a7e3-817c-9cb8-d526eedaa57c	https://acconline.austincc.edu/ultra/courses/_941496_1/grades/lti/launch?courseId=_941496_1&contentId=_28852611_1	done
56	ENGL-2311	HW	2025-06-30 00:00:00	Read: Planning Presentations/with questions (Cardon)	Read: Planning Presentations/with questions	21f40a21-a7e3-8119-81b0-ecd4a8a1e162	https://acconline.austincc.edu/ultra/courses/_941962_1/outline/lti/launch?courseId=_941962_1&contentId=_29449820_1	done
60	COSC-1301	HW	2025-07-11 00:00:00	Excel Module 04 - SAM Project 1	Do excel Module 04 - SAM Project 1	21f40a21-a7e3-81f9-9e02-fbfc128524a4	https://acconline.austincc.edu/ultra/courses/_941496_1/grades/lti/launch?courseId=_941496_1&contentId=_28852826_1	default
76	GOVT-2305	HW	2025-07-09 00:00:00	Quiz 1	will be over Chapters 1&2	22940a21-a7e3-81c6-8490-dbfeaca7ac86	https://acconline.austincc.edu/ultra/courses/_942344_1/outline/assessment/_29676236_1/overview?courseId=_942344_1	default
69	COSC-1301	Exam	2025-07-18 00:00:00	Excel 365/2021 Modules 1-4 - SAM RSA Capstone Project 1a	Do excel 365/2021 Modules 1-4 - SAM RSA Capstone Project 1a	22340a21-a7e3-8184-a0b0-ca7591247337	https://acconline.austincc.edu/ultra/courses/_941496_1/outline/lti/launch?courseId=_941496_1&contentId=_28852613_1	default
59	COSC-1301	HW	2025-07-11 00:00:00	Excel Module 03 - SAM Project 1	Do excel Module 03 - SAM Project 1	21f40a21-a7e3-8181-aa0b-e19516fce564	https://acconline.austincc.edu/ultra/courses/_941496_1/grades/lti/launch?courseId=_941496_1&contentId=_28852794_1	done
\.


--
-- Data for Name: courses; Type: TABLE DATA; Schema: public; Owner: williamfotso
--

COPY public.courses (name, code, notion_id, duration, room_number, id) FROM stdin;
Pre-Calculus	MATH-2412	17e40a21a7e380459b6fe9695d4edff9	\N	\N	1
Prog Fund III : Data Structure	COSC-2436	17f40a21a7e3801aaf8dc7bc1a0b5b3c	\N	\N	2
Fund Network	ITNW-1325	17f40a21a7e38039a931ed49406f9d9c	\N	\N	3
UNIX Op Sys I	ITSC-1307	17f40a21a7e3808d864afe449f9acec0	\N	\N	4
Info security	ITSY-1300	17e40a21a7e38097853ac873cd29f186	\N	\N	5
Technical Writing	ENGL-2311	20a40a21-a7e3-810a-959a-fdfd2cee9965	T, Th 4:00 PM - 6:00 PM	Online	12
US Governement	GOVT-2305	20a40a21-a7e3-8173-b6d6-d24599343aea	M, T, W, Th 1:30 PM - 3:30 PM	Online	13
US History I	HIST-1301	20a40a21-a7e3-815b-a682-d13d61549668	Async	Online	14
Introduction to computing	COSC-1301	20a40a21-a7e3-8112-909a-d1deff1313aa	T Th, 10am - 1pm	Online	6
\.


--
-- Data for Name: status; Type: TABLE DATA; Schema: public; Owner: williamfotso
--

COPY public.status (id, name, color) FROM stdin;
3aa77cf8-c39e-4c7b-b7d2-ab15ae43ff23	Not started	default
97903420-1e83-4b3a-9eaf-a904354c968b	In progress	blue
2fef8044-d8d7-4fcf-a3ee-393a1d558e94	Done	green
\.


--
-- Data for Name: type; Type: TABLE DATA; Schema: public; Owner: williamfotso
--

COPY public.type (id, color, name) FROM stdin;
Vn}Z	yellow	HW
oiNS	red	Exam
\.


--
-- Name: assignements_id_seq; Type: SEQUENCE SET; Schema: public; Owner: williamfotso
--

SELECT pg_catalog.setval('public.assignements_id_seq', 78, true);


--
-- Name: courses_id_seq; Type: SEQUENCE SET; Schema: public; Owner: williamfotso
--

SELECT pg_catalog.setval('public.courses_id_seq', 12, true);


--
-- Name: courses_id_seq1; Type: SEQUENCE SET; Schema: public; Owner: williamfotso
--

SELECT pg_catalog.setval('public.courses_id_seq1', 14, true);


--
-- Name: assignments assignements_pkey; Type: CONSTRAINT; Schema: public; Owner: williamfotso
--

ALTER TABLE ONLY public.assignments
    ADD CONSTRAINT assignements_pkey PRIMARY KEY (id);


--
-- Name: courses courses_code_unique; Type: CONSTRAINT; Schema: public; Owner: williamfotso
--

ALTER TABLE ONLY public.courses
    ADD CONSTRAINT courses_code_unique UNIQUE (code);


--
-- Name: courses courses_notion_id_unique; Type: CONSTRAINT; Schema: public; Owner: williamfotso
--

ALTER TABLE ONLY public.courses
    ADD CONSTRAINT courses_notion_id_unique UNIQUE (notion_id);


--
-- Name: courses courses_pkey; Type: CONSTRAINT; Schema: public; Owner: williamfotso
--

ALTER TABLE ONLY public.courses
    ADD CONSTRAINT courses_pkey PRIMARY KEY (id);


--
-- Name: assignments notion_id_unique; Type: CONSTRAINT; Schema: public; Owner: williamfotso
--

ALTER TABLE ONLY public.assignments
    ADD CONSTRAINT notion_id_unique UNIQUE (notion_id);


--
-- Name: assignments fk_course; Type: FK CONSTRAINT; Schema: public; Owner: williamfotso
--

ALTER TABLE ONLY public.assignments
    ADD CONSTRAINT fk_course FOREIGN KEY (course_code) REFERENCES public.courses(code);


--
-- PostgreSQL database dump complete
--

