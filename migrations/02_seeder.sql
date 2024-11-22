
INSERT INTO public.university (id,"name",abbreviated_name,image_url,domain_name,created_at,created_by,updated_at,updated_by,deleted_at,deleted_by) VALUES
('9a9e7fe0-37c9-40dc-bdb3-5a50adf35eaf','Universitas Multimedia Nusantara','UMN','https://assets.nsd.co.id/images/kampus/logo/umn.png','umn.ac.id','2024-11-10 16:17:33.112914+07','SYSTEM',NULL,NULL,NULL,NULL),
('a7d3c0e4-686a-45a8-861b-f3a2acc694a7','Universitas Bina Nusantara','Binus','https://kontenesia.com/wp-content/uploads/2022/07/logo-binus.png','binus.ac.id','2024-11-17 04:42:48.235505+07','SYSTEM',NULL,NULL,NULL,NULL);

INSERT INTO public."user" (id,username,email,"role","password",is_banned,is_email_verified,reputation_points,university_id,created_at,created_by,updated_at,updated_by,deleted_at,deleted_by,has_rate_university) VALUES
('06fb5492-6757-42cf-bc97-8f4569790f16','superadmin','meowhasiswa.admin@meowhasiswa.com','USER','$2a$14$PPFGM9XtN50lRK.YyKTy.eCmKzKGB3u12IeilPEhpvPx57nFP.ea2',false,false,50,NULL,'2024-11-08 16:09:11.397885+07','meowhasiswa.admin@meowhasiswa.com',NULL,NULL,NULL,NULL,false),
('f3cb8535-9977-4a1c-874a-2c0e1bed503d','inconsistence','andi.usman@student.umn.ac.id','USER','$2a$14$vvoYwuRKNiW2Fc/yfNVAzeWgPVBoQfwUUIDDqHkdsMePQ89A2sCJ.',false,false,50,'9a9e7fe0-37c9-40dc-bdb3-5a50adf35eaf','2024-11-08 12:58:21.023625+07','andi.usman@student.umn.ac.id',NULL,NULL,NULL,NULL,true),
('06f59770-6e6b-4687-9a97-64ba60ec1c68','test','test@gmail.com','USER','$2a$14$ryGhygC4SwtrXUu42CO/t.8Tgo9t5ayI9IRgdgx4S4JjL6ESpHOoe',false,false,50,'a7d3c0e4-686a-45a8-861b-f3a2acc694a7','2024-11-14 11:37:30.71319+07','test@gmail.com',NULL,NULL,NULL,NULL,false);

INSERT INTO public.subthread (id,"name",description,followers_count,image_url,university_id,is_university_subthread,created_at,created_by,updated_at,updated_by,deleted_at,deleted_by,label_color) VALUES
('c1a2f145-d7cb-4c2a-80e4-937de53697c1','LoveLife','Semua tentang cinta <3',1,'test-image-url',NULL,false,'2024-11-08 12:58:34.586771+07','SYSTEM',NULL,NULL,NULL,NULL,'#FFC0CB'),
('a6c4ff15-df98-4e65-888d-dd89efeaafd1','MarahMarah','marah marah',0,'test-image-url',NULL,false,'2024-11-10 03:14:59.48629+07','SYSTEM',NULL,NULL,NULL,NULL,'#ADFF2F'),
('7a1311cd-a25e-4270-84e7-7d37774502cf','Perkuliahan','semua tentang kuliah',1,'test-image-url',NULL,false,'2024-11-10 03:14:46.083826+07','SYSTEM',NULL,NULL,NULL,NULL,'#00FFFF'),
('f3a0b749-2e6a-48c8-ae1c-47fe63bf3497','UMN','khusus untuk mahasiswa UMN',0,'test-image-url','9a9e7fe0-37c9-40dc-bdb3-5a50adf35eaf',true,'2024-11-10 16:30:50.012861+07','andi.usman@student.umn.ac.id',NULL,NULL,NULL,NULL,'blue');

INSERT INTO public.subthread_follower (id,user_id,subthread_id,is_following,created_at,created_by,updated_at,updated_by,deleted_at,deleted_by) VALUES
('c876123e-3836-4f62-a423-bafa9f741eed','f3cb8535-9977-4a1c-874a-2c0e1bed503d','c1a2f145-d7cb-4c2a-80e4-937de53697c1',true,'2024-11-08 13:00:20.963831+07','SYSTEM',NULL,NULL,NULL,NULL),
('8034d501-7424-4b1b-a373-1bbf30ec09df','f3cb8535-9977-4a1c-874a-2c0e1bed503d','7a1311cd-a25e-4270-84e7-7d37774502cf',true,'2024-11-10 03:15:58.941941+07','SYSTEM',NULL,NULL,NULL,NULL);

INSERT INTO public.thread (id,user_id,subthread_id,title,"content",content_summary,is_active,like_count,dislike_count,comment_count,created_at,created_by,updated_at,updated_by,deleted_at,deleted_by) VALUES
('f943486f-f534-49a6-8d5d-0b456b226f18','f3cb8535-9977-4a1c-874a-2c0e1bed503d','c1a2f145-d7cb-4c2a-80e4-937de53697c1','My third love','this is a description','this is a summary',true,0,0,0,'2024-11-10 02:33:28.384301+07','andi.usman@student.umn.ac.id',NULL,NULL,NULL,NULL),
('e0b2e510-0492-479c-abb4-f1c9a044806a','f3cb8535-9977-4a1c-874a-2c0e1bed503d','7a1311cd-a25e-4270-84e7-7d37774502cf','Kuliah now','kuliah description description','this is a summary',true,0,0,0,'2024-11-10 03:16:52.744793+07','andi.usman@student.umn.ac.id',NULL,NULL,NULL,NULL),
('625468c0-0034-4554-a51f-96f8f68aa44e','f3cb8535-9977-4a1c-874a-2c0e1bed503d','c1a2f145-d7cb-4c2a-80e4-937de53697c1','My first love','this is a description','this is a summary',true,1,0,2,'2024-11-08 16:41:03.12703+07','andi.usman@student.umn.ac.id',NULL,NULL,NULL,NULL),
('a5b6fa3a-0c63-47fa-bce3-f944d55cb0a2','f3cb8535-9977-4a1c-874a-2c0e1bed503d','c1a2f145-d7cb-4c2a-80e4-937de53697c1','From appdwa','dwadwadwa','dwawdawa',true,0,0,0,'2024-11-16 17:55:27.929928+07','andi.usman@student.umn.ac.id',NULL,NULL,NULL,NULL),
('b953cea8-d0c1-4967-af18-2bbfe9d0d81b','f3cb8535-9977-4a1c-874a-2c0e1bed503d','c1a2f145-d7cb-4c2a-80e4-937de53697c1','gesgesg','ggege','gegege',true,0,0,0,'2024-11-16 18:30:28.985003+07','andi.usman@student.umn.ac.id',NULL,NULL,NULL,NULL),
('49f1a8d5-8391-40b4-9b06-be39b0d38376','f3cb8535-9977-4a1c-874a-2c0e1bed503d','7a1311cd-a25e-4270-84e7-7d37774502cf','dwadwad','dwadwad','wddw',true,0,0,0,'2024-11-16 18:46:15.96025+07','andi.usman@student.umn.ac.id',NULL,NULL,NULL,NULL),
('b1887735-1d89-48bf-8289-a60454bd8154','f3cb8535-9977-4a1c-874a-2c0e1bed503d','7a1311cd-a25e-4270-84e7-7d37774502cf','From app','wdadwa','wdadaw',true,0,0,0,'2024-11-16 23:33:58.500297+07','andi.usman@student.umn.ac.id',NULL,NULL,NULL,NULL),
('3940cf54-4ae3-47ee-b8e6-e344206e6d9a','06f59770-6e6b-4687-9a97-64ba60ec1c68','c1a2f145-d7cb-4c2a-80e4-937de53697c1','My second love','this is a description','this is a summary',true,0,0,0,'2024-11-10 02:33:24.5118+07','andi.usman@student.umn.ac.id',NULL,NULL,NULL,NULL),
('ada7d3c3-8c34-468e-9fc4-0ed99286285f','f3cb8535-9977-4a1c-874a-2c0e1bed503d','7a1311cd-a25e-4270-84e7-7d37774502cf','From app ','odiwnaodwandw','adwadwawa',true,0,0,0,'2024-11-17 20:24:40.991267+07','andi.usman@student.umn.ac.id',NULL,NULL,NULL,NULL);

INSERT INTO public.thread_activity (id,actor_id,actor_email,actor_username,thread_id,"action",created_at,created_by) VALUES
('dc50b7da-0de7-479e-b06e-8cf2a3766a4e','f3cb8535-9977-4a1c-874a-2c0e1bed503d','andi.usman@student.umn.ac.id','inconsistence','625468c0-0034-4554-a51f-96f8f68aa44e','LIKE','2024-11-10 19:58:47.654702+07','andi.usman@student.umn.ac.id'),
('fcf1d117-ec4c-498a-8a30-9c242c513c1a','f3cb8535-9977-4a1c-874a-2c0e1bed503d','andi.usman@student.umn.ac.id','inconsistence','625468c0-0034-4554-a51f-96f8f68aa44e','UNLIKE','2024-11-10 19:59:00.907747+07','andi.usman@student.umn.ac.id'),
('ba05366a-eab8-4a5e-bbe9-c6fd0ccf7220','f3cb8535-9977-4a1c-874a-2c0e1bed503d','andi.usman@student.umn.ac.id','inconsistence','625468c0-0034-4554-a51f-96f8f68aa44e','DISLIKE','2024-11-10 20:14:27.734713+07','andi.usman@student.umn.ac.id'),
('3d2ad90e-d9b9-4656-a2c6-4f36dae6dc12','f3cb8535-9977-4a1c-874a-2c0e1bed503d','andi.usman@student.umn.ac.id','inconsistence','625468c0-0034-4554-a51f-96f8f68aa44e','UNDISLIKE','2024-11-10 20:14:59.88312+07','andi.usman@student.umn.ac.id'),
('9e877f5e-78f4-471b-bc76-4f86e5117c80','f3cb8535-9977-4a1c-874a-2c0e1bed503d','andi.usman@student.umn.ac.id','inconsistence','625468c0-0034-4554-a51f-96f8f68aa44e','DISLIKE','2024-11-10 20:15:08.026527+07','andi.usman@student.umn.ac.id'),
('9cbfa418-c3f8-4f81-8a49-364e0e354f88','f3cb8535-9977-4a1c-874a-2c0e1bed503d','andi.usman@student.umn.ac.id','inconsistence','625468c0-0034-4554-a51f-96f8f68aa44e','LIKE','2024-11-10 20:15:26.546901+07','andi.usman@student.umn.ac.id');

INSERT INTO public.thread_comment (id,user_id,thread_id,"content",like_count,dislike_count,reply_count,created_at,created_by,updated_at,updated_by,deleted_at,deleted_by) VALUES
('8dff604c-33a5-4c23-a953-f77fb81eafd6','f3cb8535-9977-4a1c-874a-2c0e1bed503d','625468c0-0034-4554-a51f-96f8f68aa44e','this is a comment',0,0,1,'2024-11-15 04:15:52.091383+07','andi.usman@student.umn.ac.id',NULL,NULL,NULL,NULL);

INSERT INTO public.thread_comment_reply (id,user_id,thread_id,thread_comment_id,"content",like_count,dislike_count,created_at,created_by,updated_at,updated_by,deleted_at,deleted_by) VALUES
('91f64b2f-9b08-4e0c-ac72-44b1aa1b624a','f3cb8535-9977-4a1c-874a-2c0e1bed503d','625468c0-0034-4554-a51f-96f8f68aa44e','8dff604c-33a5-4c23-a953-f77fb81eafd6','this is a comment',0,0,'2024-11-15 04:16:43.739085+07','andi.usman@student.umn.ac.id',NULL,NULL,NULL,NULL);

INSERT INTO public.university_rating (id,user_id,university_id,title,"content",university_major,facility_rating,student_organization_rating,social_environment_rating,education_quality_rating,price_to_value_rating,overall_rating,created_at,created_by,updated_at,updated_by,deleted_at,deleted_by) VALUES
('3ed7da5e-a3b1-4c73-a55d-c06418b8b01d','f3cb8535-9977-4a1c-874a-2c0e1bed503d','9a9e7fe0-37c9-40dc-bdb3-5a50adf35eaf','Great University Experience','This university provides an excellent balance between academics and social life. The facilities are top-notch, and there are numerous student organizations to join.','Informatika',5,4,5,4,3,4.20,'2024-11-17 04:41:28.348082+07','andi.usman@student.umn.ac.id','2024-11-17 04:41:28.348082+07','andi.usman@student.umn.ac.id',NULL,NULL),
('4c9de885-e2f1-4810-998e-a6c6f5aa56ed','06f59770-6e6b-4687-9a97-64ba60ec1c68','a7d3c0e4-686a-45a8-861b-f3a2acc694a7','Very Nice','This university provides an excellent balance between academics and social life. The facilities are top-notch, and there are numerous student organizations to join.','Akuntansi',5,5,5,5,5,5.00,'2024-11-17 04:45:11.344989+07','andi.usman@student.umn.ac.id','2024-11-17 04:45:11.344989+07','andi.usman@student.umn.ac.id',NULL,NULL);

INSERT INTO public.university_rating_point (id,university_rating_id,"type","content",created_at,created_by,updated_at,updated_by) VALUES
('e9b546f7-4c45-4c0c-8e2e-85db3d382907','3ed7da5e-a3b1-4c73-a55d-c06418b8b01d','PRO','Modern facilities','2024-11-17 04:41:28.348082+07','andi.usman@student.umn.ac.id','2024-11-17 04:41:28.348082+07','andi.usman@student.umn.ac.id'),
('8142e2b2-3cfc-4f82-b125-fabf8b0d246d','3ed7da5e-a3b1-4c73-a55d-c06418b8b01d','PRO','Active student organizations','2024-11-17 04:41:28.348082+07','andi.usman@student.umn.ac.id','2024-11-17 04:41:28.348082+07','andi.usman@student.umn.ac.id'),
('9f500806-dbbb-4031-8cc7-f69025ff0816','3ed7da5e-a3b1-4c73-a55d-c06418b8b01d','PRO','Diverse social environment','2024-11-17 04:41:28.348082+07','andi.usman@student.umn.ac.id','2024-11-17 04:41:28.348082+07','andi.usman@student.umn.ac.id'),
('04d6b8bd-3147-4499-af25-9091d10d64f4','3ed7da5e-a3b1-4c73-a55d-c06418b8b01d','CON','High tuition fees','2024-11-17 04:41:28.348082+07','andi.usman@student.umn.ac.id','2024-11-17 04:41:28.348082+07','andi.usman@student.umn.ac.id'),
('caa23616-609c-41d3-849a-ce56891566f5','3ed7da5e-a3b1-4c73-a55d-c06418b8b01d','CON','Limited parking spaces','2024-11-17 04:41:28.348082+07','andi.usman@student.umn.ac.id','2024-11-17 04:41:28.348082+07','andi.usman@student.umn.ac.id'),
('19fd82e6-1bbd-46b3-b18e-d982831383d0','4c9de885-e2f1-4810-998e-a6c6f5aa56ed','PRO','Good vibes','2024-11-17 04:46:15.898224+07','andi.usman@student.umn.ac.id','2024-11-17 04:46:15.898224+07','andi.usman@student.umn.ac.id'),
('c42c6327-baec-4dbe-88a4-1a6e95d3c174','4c9de885-e2f1-4810-998e-a6c6f5aa56ed','CON','Mahal','2024-11-17 04:46:15.898224+07','andi.usman@student.umn.ac.id','2024-11-17 04:46:15.898224+07','andi.usman@student.umn.ac.id');
