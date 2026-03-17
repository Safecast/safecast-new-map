-- Seed translations from translations.json
-- Generated automatically - run after add_translations_table.sql

BEGIN;

INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_archive_desc', 'تنزّل حزمة tgz تضم كل ملفات .json المنشورة عندما يكون أرشيف JSON مفعلاً.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_archive_link', 'تنزيل الأرشيف الأسبوعي') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_archive_note', 'إذا كان الأرشيف معطلاً سيُرجع الخادم الخطأ HTTP 404 غير موجود.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_archive_title', 'حزمة الأرشيف الأسبوعية') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_latest_desc', 'يعيد أحدث العلامات قرب خط العرض، خط الطول ونصف القطر المطلوب بالأمتار.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_latest_link', 'الأحدث قرب طوكيو') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_latest_note', 'عدّل قيم lat وlon وradius_m للتركيز على منطقتك.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_latest_title', 'أحدث القياسات القريبة') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_root_desc', 'يسرد البيانات الوصفية وإحصاءات مجموعة البيانات والروابط إلى كل واجهة أخرى.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_root_link', 'افتح /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_root_note', 'ابدأ من هنا لاستكشاف المجموعات وحالة الخادم.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_root_title', 'فهرس قابل للقراءة آليًا') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_track_index_desc', 'يجلب المسار رقم N ويعيد نفس مستند JSON كما في /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_track_index_link', 'جلب الفهرس 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_track_index_note', 'استبدل الرقم لجلب إدخال مختلف. تُبث الاستجابات كسطور JSON متتابعة.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_track_index_title', 'تنزيل مسار حسب الفهرس') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_tracks_desc', 'يبث فهرسًا للمسارات المنشورة يتضمن الأسماء والأعداد وروابط التنزيل.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_tracks_link', 'عرض أول المسارات') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_tracks_note', 'استخدم معاملي limit وoffset للتصفح عبر القوائم الطويلة.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_example_tracks_title', 'ملخصات كل المسارات') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_examples_heading', 'نقاط نهاية مميزة') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_examples_note', 'تستجيب جميع الروابط بصيغة JSON. قد يقوم المتصفح بتنزيل ملفات كبيرة عندما تبث الواجهات البيانات.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_intro', 'تعكس واجهة JSON البيانات المعروضة على الخريطة. يفتح كل مثال في علامة تبويب جديدة لتتمكن من معاينة الاستجابة الخام.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_more_docs', 'تحتاج إلى توثيق أعمق؟') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_more_docs_link_label', 'افتح الدليل المفصل') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'api_title', 'بدء سريع لواجهة API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'attr_api', 'واجهة API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'attr_legal', 'معلومات قانونية') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'attr_license', 'الترخيص') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'attr_sources', 'مصادر البيانات') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'back_to_all_tracks', 'العودة إلى خريطة المسارات الموحّدة.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'count_rate', 'معدل العدّ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'date_slider_tooltip', 'رشّح القياسات حسب التاريخ أو نطاق السنة.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'date_time', 'التاريخ والوقت') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'description', 'خريطة الإشعاع لتشيتشا') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'download_track_cim', 'تحميل المسار (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'duration_days', '[[count]] يوم') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'duration_hours', '[[count]] س') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'duration_minutes', '[[count]] د') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'duration_months', '[[count]] شهر') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'duration_weeks', '[[count]] أسبوع') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'error_during_upload', 'خطأ أثناء الرفع!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'error_processing_files', 'خطأ في معالجة الملفات!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'files_uploaded', 'تم رفع الملفات') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'geolocation_not_supported', 'متصفحك لا يدعم تحديد الموقع الجغرافي.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'github_link_tooltip_desc', 'مشروع مفتوح المصدر يديره المجتمع.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'github_link_tooltip_title', 'خريطة تشيتشا للنظائر على GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'github_link_tooltip_version', 'الإصدار الحالي: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'home_explore_global', 'تخطّي واستكشاف الخريطة العالمية') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'home_or', 'أو') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'home_search_placeholder', 'أدخل مدينة أو منطقة أو بلد...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'home_subtitle', 'أدخل موقعاً للبدء.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'home_title', 'خريطة Safecast الإشعاعية') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'home_use_location', 'استخدم موقعي') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'legal_contact', 'لإرسال الملاحظات راسلونا على:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'legal_full', 'التمهيد. نبني خريطة مفتوحة يتقاسم فيها الناس من كل أنحاء العالم قراءات مقاييس الجرعات خدمةً للصالح العام — للعلم والبيئة والتعليم والسلامة. بمشاركتك لبياناتك تساعد الكثيرين. نرجو أن تتعامل مع هذا العمل المشترك بمحبة ومسؤولية.

1) المسؤولية. أنت مسؤول عن دقة المحتوى الذي ترسله. تُنشر البيانات وتُستخدم «كما هي». لا تؤكد الخدمة صحتها أو اكتمالها أو ملاءمتها لأي غرض، ولا تتحمل أية تبعات قد تنجم عن استخدامها.

2) الانفتاح والترخيص. حين تشارك القياسات أو التواريخ أو المواقع التقريبية أو طراز الجهاز أو أي معلومات أخرى، فأنت تدرك أنها تصبح متاحة للجميع ويمكن استخدامها بحرية بموجب ترخيص البيانات CC0 1.0 (الملك العام)، بينما يبقى الكود متاحًا برخصة MIT. يبقى حقك في النسبة محفوظًا؛ لا يقدم مقابل مالي؛ وإعادة التوزيع اللاحقة من قبل أطراف ثالثة خارج نطاق سيطرتنا.

3) «كما هي» ومن دون تحقق. تُنشر المواد من دون مراجعة مسبقة. لا يمكننا ضمان معايرة الأجهزة أو خلوها من الأخطاء. تُقدَّم المعلومات لأغراض بحثية ولا تُعد توصية مهنية.

4) الخصوصية والإشراف. حفاظًا على الأمان والثقة قد نعمّم الأوقات والإحداثيات، وقد نحذف أو نُخفي البيانات الوصفية التقنية. ويجوز لنا إخفاء أو حذف المواد التي نراها بعقلانية بريدًا عشوائيًا أو مزيفة أو مخالفة للقانون أو معطلة للخدمة. نتعامل بعناية مع القياسات الموثوقة ونسعى إلى حفظها.

5) ملفات تعريف الارتباط. يستخدم الموقع ملف تعريف ارتباط تقنيًا قصير الأجل للجلسة؛ يختفي عند انتهاء زيارتك. لا نحتفظ بأي آثار أخرى.

أصدقاؤنا، هذه الخريطة ثمرة جهد جماعي وقلوب منفتحة. عامِلوها كمسودة للتضاريس لا كخريطة دقيقة. إن لامسكم هذا العمل، فانضموا إلينا — معًا نستطيع أن نجعله أفضل.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'legal_title', 'معلومات قانونية') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'legend_attention', 'انتباه') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'legend_button_tooltip', 'افتح دليل مستويات الإشعاع.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'legend_danger', 'خطر') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'legend_full_ar', 'هذا المقياس يوضح مدى أمان المكان للحياة والماء والطعام.
تذكّر: قد تكون القياسات ناقصة؛ استخدمها كدليل فقط.

أخضر (0–11 µR/h)
خلفية طبيعية.
• ماء الآبار غالبًا آمن.
• يمكن الزراعة دون فحص.

أصفر (11–30 µR/h)
خلفية مرتفعة.
• افحص الماء والتربة.
• اختبر أي طعام قبل تناوله.

أحمر (30–100 µR/h)
تلوث خطير.
• لا تشرب الماء.
• الزراعة أو الأكل هنا خطر؛ الفحوص المخبرية ضرورية.

أسود (>100 µR/h)
منطقة حرجة.
• الماء والطعام غير صالحين.
• الإقامة الطويلة ممنوعة؛ الزيارة القصيرة مع حماية فقط.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'legend_safe', 'آمن') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'legend_title', 'مفتاح') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'license_full', 'ينمو هذا المشروع تحت رخصة <a href="/LICENSE" target="_blank">MIT License</a>. يوجد النص الكامل في جذر المستودع وعلى موقع MIT. يمكنك دراسة الكود ومشاركته وتعديله ما دامت هذه الحريات ترافق عملك. تُنشر بيانات الأبحاث بموجب ترخيص Creative Commons 1.0 لضمان بقاء القياسات في الملك العام.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'license_title', 'الترخيص') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_chart_all', 'جميع القراءات') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_chart_averaged', 'متوسط خلال [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_chart_close', 'إغلاق') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_chart_day', 'آخر 24 ساعة') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_chart_month', 'آخر 30 يوماً') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_chart_link', 'فتح رسوم الإشعاع البيانية') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_country', 'البلد') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_desc_generic', 'مستشعر Safecast هذا') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_desc_location', 'يرسل بيانات من [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_desc_location_unknown', 'يرسل بيانات من منطقة غير معروفة') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_desc_named', 'الجهاز [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_desc_no_tube', 'يراقب مستويات الإشعاع.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_desc_radiation_sensor', 'هذا مستشعر إشعاع.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_desc_transport_air', 'أثناء الطيران') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_desc_transport_bike', 'أثناء التنقل بالدراجة') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_desc_transport_car', 'أثناء التنقل بالسيارة') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_desc_transport_unknown', 'أثناء البقاء في المكان') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_desc_transport_walk', 'أثناء المشي') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_desc_tube', 'باستخدام كاشف [[tube]].') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_detector', 'الكاشف') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_device_id', 'معرّف الجهاز') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_device_label', 'الجهاز') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_extra_intro', 'البيئة') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_history_error', 'تعذّر تحميل السجل الآن.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_humidity', 'الرطوبة') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_last_seen', 'آخر قراءة') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_local_time', 'التوقيت المحلي الآن') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_no_data', 'لا توجد بيانات مسجلة في هذه الفترة.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_pressure', 'ضغط الهواء') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_temperature', 'درجة حرارة الهواء') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_temperature_f', 'درجة حرارة الهواء (°ف)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_title', 'مستشعر Safecast الفوري') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_transport', 'الحركة') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_transport_air', 'طائرة') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_transport_bike', 'دراجة أو سكوتر') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_transport_car', 'سيارة أو شاحنة صغيرة') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_transport_unknown', 'غير محدد') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'live_marker_transport_walk', 'مشياً على الأقدام') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'locate_button_tooltip', 'مركِز الخريطة على موقعي') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'location_error', 'حدث خطأ أثناء الحصول على الموقع.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'location_permission_denied', 'تم رفض الوصول إلى الموقع.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'location_timeout', 'انتهت مهلة طلب الموقع.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'location_unavailable', 'الموقع غير متاح.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'no_results_found', 'لم يتم العثور على نتائج') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'processing_complete', 'اكتملت المعالجة!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'processing_on_server', 'جارٍ المعالجة على الخادم...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'qr_button_tooltip', 'رمز QR لرابط هذه المنطقة من الخريطة.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'radiation_dose', 'معدل الجرعة') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'search_error', 'خطأ في البحث. يرجى المحاولة مرة أخرى.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'search_input_tooltip', 'ابحث عن مكان بكتابة الأحرف الأولى من اسمه. ستظهر قائمة بالاقتراحات.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'select_files', 'يرجى اختيار ملف واحد على الأقل') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'short_link_tooltip', 'انقر لنسخ رابط قصير للمشاركة') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'sources_full', 'نحن نشكر كل من يشارك القياسات.

عمليات الرفع المجهولة ترسم مسارات هادئة على الخريطة.
<a href="https://safecast.org" target="_blank">Safecast</a> ترعى أرشيفًا عالميًا للقراءات.
<a href="https://atomfast.net" target="_blank">Atomfast</a> يبقي Atomcloud متقدًا.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> يجمع رؤى Radiacode.

كل مساهمة توسّع الصورة المشتركة؛ ندعوك بحرارة لإضافة مساهمتك.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'sources_title', 'مصادر البيانات') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'speed', 'السرعة') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'speed_filter_tooltip_accuracy', 'القياسات الأبطأ تبقى أقرب إلى الأرض، لذا بيانات المشاة هي الأدق.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'speed_filter_tooltip_car', 'سيارة: من 7 إلى 200 كم/س للقيادة والمسوح المتنقلة.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'speed_filter_tooltip_intro', 'اختر القياسات التي تظهر بحسب سرعة الحركة.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'speed_filter_tooltip_live', 'قلب Safecast: بيانات فورية من safecast.org. بدّل الخيار لإظهار القراءات الحية أو إخفائها.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'speed_filter_tooltip_ped', 'مشاة: أقل من 7 كم/س للقياسات أثناء المشي أو التوقف.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'speed_filter_tooltip_plane', 'طائرة: +200 كم/س لعمليات المسح الجوية.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'speed_filter_tooltip_title', 'مرشحات السرعة') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'theme_toggle_tooltip', 'بدّل بين سمة الخريطة الفاتحة والداكنة.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'title', 'خريطة نظائر تشيتشا — خريطة إشعاعية لكوكب الأرض') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'track_id', 'معرّف المسار') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'upload_button', 'رفع [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'upload_button_tooltip', 'أضف مسار القياسات الخاص بك إلى الخريطة. الصيغ المدعومة: .kml، .kmz، .gpx، .csv، .rctrk، .json، .log. يمكنك رفع عدة ملفات، وبعد الرفع ستفتح صفحة المسار.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'upload_error', 'خطأ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'waiting_for_server', 'بانتظار الخادم...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ar', 'your_location', 'موقعك') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_archive_desc', 'Stáhne archiv tgz se všemi zveřejněnými soubory .json, pokud je JSON archiv zapnutý.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_archive_link', 'Stáhnout týdenní archiv') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_archive_note', 'Pokud je archiv vypnutý, server vrátí chybu HTTP 404 Nenalezeno.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_archive_title', 'Týdenní archivní balíček') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_latest_desc', 'Vrátí nejnovější body v okolí zadané zeměpisné šířky, délky a poloměru v metrech.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_latest_link', 'Nejnovější u Tokia') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_latest_note', 'Upravte hodnoty lat, lon a radius_m pro vlastní oblast.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_latest_title', 'Nejnovější měření v okolí') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_root_desc', 'Vypisuje metadata, statistiky datasetu a odkazy na všechny ostatní koncové body.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_root_link', 'Otevřít /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_root_note', 'Začněte zde a zjistěte kolekce i stav serveru.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_root_title', 'Strojově čitelný index') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_track_index_desc', 'Načte N-tou trasu a vrátí stejný JSON jako /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_track_index_link', 'Získat index 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_track_index_note', 'Nahraďte číslo pro jinou položku. Odpověď streamuje jako JSON po řádcích.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_track_index_title', 'Stáhnout trasu podle pořadí') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_tracks_desc', 'Streamuje katalog zveřejněných tras včetně názvů, počtů a odkazů ke stažení.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_tracks_link', 'Zobrazit první trasy') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_tracks_note', 'Použijte parametry limit a offset pro stránkování dlouhých seznamů.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_example_tracks_title', 'Souhrny všech tras') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_examples_heading', 'Doporučené koncové body') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_examples_note', 'Všechny odkazy vracejí JSON. Při streamování mohou prohlížeče stahovat velké soubory.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_intro', 'JSON API zrcadlí data zobrazená na mapě. Každý příklad se otevře na nové kartě, abyste mohli prohlédnout syrovou odpověď.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_more_docs', 'Potřebujete podrobnější informace?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_more_docs_link_label', 'Otevřít podrobnou příručku') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'api_title', 'Rychlý start API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'attr_legal', 'Právní informace') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'attr_license', 'Licence') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'attr_sources', 'Zdroje dat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'back_to_all_tracks', 'Zpět na kombinovanou mapu tras.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'count_rate', 'Četnost impulzů') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'date_slider_tooltip', 'Filtrovat měření podle data nebo roku.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'date_time', 'Datum a čas') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'description', 'Safecast Nová Mapa') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'download_track_cim', 'Stáhnout trasu (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'duration_days', '[[count]] dní') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'duration_hours', '[[count]] hod') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'duration_minutes', '[[count]] min') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'duration_months', '[[count]] měsíců') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'duration_weeks', '[[count]] týdnů') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'error_during_upload', 'Chyba při nahrávání!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'error_processing_files', 'Chyba při zpracování souborů!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'files_uploaded', 'Soubory nahrány') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'geolocation_not_supported', 'Váš prohlížeč nepodporuje geolokaci.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'github_link_tooltip_desc', 'Projekt s otevřeným zdrojovým kódem spravovaný komunitou.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'github_link_tooltip_title', 'Safecast New Map na GitHubu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'github_link_tooltip_version', 'Aktuální verze: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'home_explore_global', 'Přeskočit a prozkoumat celou mapu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'home_or', 'nebo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'home_search_placeholder', 'Zadejte město, region nebo zemi...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'home_subtitle', 'Zadejte místo pro začátek.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'home_title', 'Safecast radiologická mapa') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'home_use_location', 'Použít moji polohu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'legal_contact', 'Máte-li zpětnou vazbu, napište na:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'legal_full', 'Preambule. Společně tvoříme otevřenou mapu, na níž lidé z celého světa sdílejí dozimetrická měření pro dobro všech – pro vědu, životní prostředí, vzdělávání i bezpečnost. Když zveřejníte svá data, pomáháte mnoha lidem. Prosíme, pečujte o toto společné dílo.

1) Odpovědnost. Za správnost a obsah zasílaných informací nesete odpovědnost vy sami. Data se zveřejňují a používají „tak, jak jsou“. Služba nepotvrzuje ani negarantuje jejich správnost, úplnost či vhodnost pro jakýkoli účel a neodpovídá za následky jejich využití.

2) Otevřenost a licence. Sdílením měření, dat, přibližné polohy, modelu zařízení nebo dalších údajů berete na vědomí, že budou dostupné všem a lze je svobodně používat podle licence CC0 1.0 (Public Domain). Kód zůstává pod licencí MIT. Autorství vám zůstává; odměna se neposkytuje; další šíření třetími stranami nemáme pod kontrolou.

3) „Jak jsou“ a bez ověření. Příspěvky se zveřejňují bez předchozí kontroly. Nemůžeme zaručit kalibraci přístrojů ani bezchybnost údajů. Informace slouží k výzkumným účelům a nepředstavují odborné doporučení.

4) Soukromí a moderace. Kvůli bezpečnosti a důvěře můžeme časy a souřadnice zobecňovat a technická metadata odstraňovat či anonymizovat. Materiály, které jsou podle našeho rozumného uvážení spam, padělek, porušují zákon nebo narušují službu, můžeme skrýt či smazat. Poctivých měření si vážíme a usilujeme o jejich zachování.

5) Cookies. Web používá pouze krátkodobé technické session cookie; po skončení návštěvy zmizí. Žádné další stopy neukládáme.

Přátelé, tato mapa je výsledkem společného úsilí a otevřených srdcí. Vnímejte ji jako náčrt krajiny, nikoli jako přesný plán. Pokud vám naše práce dává smysl, přidejte se – společně ji vylepšíme.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'legal_title', 'Právní informace') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'legend_attention', 'Pozor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'legend_button_tooltip', 'Otevřít legendu úrovní radiace.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'legend_danger', 'Nebezpečí') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'legend_full_cs', 'Tato stupnice ukazuje, jak bezpečné je dané místo pro lidi, vodu a potraviny.
Pamatujte: měření nemusí být úplné; používejte jej pouze jako vodítko.

Zelená (0–11 µR/h)
Téměř přirozené pozadí.
• Voda ze studní je obvykle bezpečná.
• Rostliny lze pěstovat bez dalších kontrol.

Žlutá (11–30 µR/h)
Zvýšené pozadí; buďte obezřetní.
• Zkontrolujte vodu a půdu.
• Zeleninu, houby a jiné plodiny před konzumací otestujte.

Červená (30–100 µR/h)
Vážná kontaminace.
• Nepijte místní vodu.
• Pěstování a konzumace místních potravin je rizikové; laboratorní testy jsou nezbytné.

Černá (>100 µR/h)
Kritická zóna.
• Voda ani potraviny nejsou použitelné.
• Dlouhodobý pobyt nepřichází v úvahu; pouze krátké návštěvy s ochranou.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'legend_safe', 'Bezpečné') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'legend_title', 'Legenda') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'license_full', 'Tento projekt funguje pod licencí <a href="/LICENSE" target="_blank">MIT License</a>. Celý text najdete v kořeni repozitáře a na webu MIT. Kód můžete studovat, sdílet i upravovat, pokud si tyto svobody ponechají i vaše úpravy. Výzkumná data zveřejňujeme pod licencí Creative Commons 1.0, takže měření zůstávají ve veřejném prostoru.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'license_title', 'Licence') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_chart_all', 'Všechna měření') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_chart_averaged', 'Průměr za [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_chart_close', 'Zavřít') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_chart_day', 'Posledních 24 hodin') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_chart_month', 'Posledních 30 dní') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_chart_link', 'Otevřít grafy radiace') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_country', 'Země') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_desc_generic', 'Tento senzor Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_desc_location', 'měří z [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_desc_location_unknown', 'měří z neznámé oblasti') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_desc_named', 'Zařízení [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_desc_no_tube', 'sleduje úroveň radiace.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_desc_radiation_sensor', 'Toto je senzor radiace.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_desc_transport_air', 'za letu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_desc_transport_bike', 'na kole') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_desc_transport_car', 'v autě') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_desc_transport_unknown', 'na místě') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_desc_transport_walk', 'pěšky') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_desc_tube', 's detektorem [[tube]].') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_detector', 'Detektor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_device_id', 'ID zařízení') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_device_label', 'Zařízení') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_extra_intro', 'Prostředí') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_history_error', 'Historii se nyní nepodařilo načíst.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_humidity', 'Vlhkost') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_last_seen', 'Poslední měření') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_local_time', 'Místní čas') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_no_data', 'V tomto období nebyla zaznamenána žádná data.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_pressure', 'Atmosférický tlak') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_temperature', 'Teplota vzduchu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_temperature_f', 'Teplota vzduchu (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_title', 'Safecast senzor v reálném čase') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_transport', 'Pohyb') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_transport_air', 'Letadlo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_transport_bike', 'Kolo nebo koloběžka') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_transport_car', 'Auto nebo dodávka') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_transport_unknown', 'Nezjištěno') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'live_marker_transport_walk', 'Pěšky') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'locate_button_tooltip', 'Vycentrovat mapu na moji polohu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'location_error', 'Při získávání polohy došlo k chybě.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'location_permission_denied', 'Přístup k poloze byl odepřen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'location_timeout', 'Časový limit požadavku na polohu vypršel.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'location_unavailable', 'Poloha není k dispozici.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'no_results_found', 'Nebyly nalezeny žádné výsledky') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'processing_complete', 'Zpracování dokončeno!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'processing_on_server', 'Zpracování na serveru...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'qr_button_tooltip', 'QR kód odkazu na tuto část mapy.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'radiation_dose', 'Dávkový příkon') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'search_error', 'Chyba vyhledávání. Zkuste to prosím znovu.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'search_input_tooltip', 'Vyhledejte místo zadáním prvních písmen jeho názvu. Zobrazí se seznam návrhů.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'select_files', 'Vyberte prosím alespoň jeden soubor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'short_link_tooltip', 'Kliknutím zkopírujete krátký sdílený odkaz') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'sources_full', 'Všem, kdo sdílejí měření, děkujeme.

Anonymní nahrávky kreslí na mapě tiché stopy.
<a href="https://safecast.org" target="_blank">Safecast</a> pečuje o celosvětový archiv měření.
<a href="https://atomfast.net" target="_blank">Atomfast</a> udržuje službu Atomcloud v chodu.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> shromažďuje poznatky z Radiacode.

Každý příspěvek rozšiřuje společný obraz; budeme rádi, když přidáte ten svůj.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'sources_title', 'Zdroje dat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'speed', 'Rychlost') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'speed_filter_tooltip_accuracy', 'Pomalejší měření zůstávají nejblíže zemi, proto jsou pěší data nejpřesnější.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'speed_filter_tooltip_car', 'Auto: 7–200 km/h pro jízdy a mobilní průzkumy.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'speed_filter_tooltip_intro', 'Vyberte, která měření se mají zobrazit podle cestovní rychlosti.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'speed_filter_tooltip_live', 'Srdce Safecast: data v reálném čase ze safecast.org. Přepínač zobrazuje nebo skrývá živá měření.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'speed_filter_tooltip_ped', 'Pěší: pod 7 km/h pro pěší nebo stojící měření.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'speed_filter_tooltip_plane', 'Letadlo: 200+ km/h pro letecké průzkumy.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'speed_filter_tooltip_title', 'Filtry rychlosti') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'theme_toggle_tooltip', 'Přepínejte mezi světlým a tmavým vzhledem mapy.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'title', 'Safecast Nová Mapa — Radiologická mapa planety Země') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'track_id', 'ID trasy') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'upload_button', 'Nahrát [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'upload_button_tooltip', 'Přidejte svou měřicí trasu do mapy. Podporované formáty: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Můžete nahrát více souborů a po nahrání se otevře stránka trasy.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'upload_error', 'Chyba') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'waiting_for_server', 'Čekání na server...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('cs', 'your_location', 'Vaše poloha') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_archive_desc', 'Downloader en tgz-pakke med alle offentliggjorte .json-filer når JSON-arkivet er aktiveret.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_archive_link', 'Download ugentligt arkiv') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_archive_note', 'Hvis arkivet er deaktiveret returnerer serveren HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_archive_title', 'Ugentlig arkivpakke') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_latest_desc', 'Returnerer de nyeste punkter nær den angivne breddegrad, længdegrad og radius i meter.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_latest_link', 'Seneste nær Tokyo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_latest_note', 'Juster lat, lon og radius_m for dit eget område.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_latest_title', 'Seneste målinger i nærheden') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_root_desc', 'Viser metadata, datasætstatistikker og links til alle andre endepunkter.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_root_link', 'Åbn /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_root_note', 'Start her for at finde samlinger og serverstatus.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_root_title', 'Maskinlæsbart indeks') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_track_index_desc', 'Finder det N''te spor og returnerer det samme JSON-dokument som /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_track_index_link', 'Hent indeks 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_track_index_note', 'Erstat tallet for at hente en anden post. Svar streames som JSON-linjer.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_track_index_title', 'Download spor efter indeks') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_tracks_desc', 'Streamer et katalog over offentliggjorte spor med navne, antal og download-links.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_tracks_link', 'Vis første spor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_tracks_note', 'Brug parametrene limit og offset til at paginere lange lister.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_example_tracks_title', 'Alle sporoversigter') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_examples_heading', 'Udvalgte endepunkter') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_examples_note', 'Alle links svarer med JSON. Når endepunkter streamer data, kan browseren downloade store filer.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_intro', 'JSON-API''et afspejler dataene på kortet. Hvert eksempel åbnes i en ny fane, så du kan inspicere det rå svar.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_more_docs', 'Brug for mere detaljeret information?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_more_docs_link_label', 'Åbn den detaljerede vejledning') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'api_title', 'API hurtig start') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'attr_legal', 'Juridiske oplysninger') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'attr_license', 'Licens') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'attr_sources', 'Datakilder') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'back_to_all_tracks', 'Tilbage til det kombinerede sporkort.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'count_rate', 'Tællehastighed') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'date_slider_tooltip', 'Filtrér målinger efter dato eller år.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'date_time', 'Dato og tid') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'description', 'Safecasts strålingskort') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'download_track_cim', 'Hent spor (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'duration_days', '[[count]] dage') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'duration_hours', '[[count]] t') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'duration_minutes', '[[count]] min') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'duration_months', '[[count]] måneder') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'duration_weeks', '[[count]] uger') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'error_during_upload', 'Fejl under upload!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'error_processing_files', 'Fejl under behandling af filer!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'files_uploaded', 'Filer uploadet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'geolocation_not_supported', 'Din browser understøtter ikke geolokation.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'github_link_tooltip_desc', 'Open source-projekt vedligeholdt af fællesskabet.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'github_link_tooltip_title', 'Safecast New Map på GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'github_link_tooltip_version', 'Nuværende version: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'home_explore_global', 'Spring over og udforsk verdenskortet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'home_or', 'eller') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'home_search_placeholder', 'Indtast en by, region eller land...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'home_subtitle', 'Indtast en placering for at begynde.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'home_title', 'Safecast radiologisk kort') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'home_use_location', 'Brug min placering') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'legal_contact', 'Har du feedback, så skriv til:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'legal_full', 'Forord. Vi bygger et åbent kort, hvor mennesker verden over deler dosimetermålinger til gavn for fællesskabet – for videnskab, miljø, uddannelse og sikkerhed. Når du offentliggør dine data, gør du en forskel for mange. Pas godt på dette fælles projekt.

1) Ansvar. Du står selv inde for nøjagtighed og indhold i de oplysninger, du indsender. Data offentliggøres og bruges “som de er”. Tjenesten bekræfter eller garanterer ikke deres korrekthed, fuldstændighed eller egnethed til noget formål og påtager sig ikke ansvar for følgerne af deres brug.

2) Åbenhed og licens. Når du deler målinger, datoer, omtrentlige positioner, modelnavne eller andre oplysninger, accepterer du, at de bliver tilgængelige for alle og kan bruges frit under CC0 1.0-licensen (Public Domain) for data. Koden udgives fortsat under MIT-licensen. Du beholder ophavsretten; der udbetales ingen betaling; videre distribution fra tredjeparter er uden for vores kontrol.

3) “Som de er” og uden forhåndstjek. Materiale udgives uden forudgående gennemgang. Vi kan ikke garantere kalibrering af udstyr eller fravær af fejl. Oplysningerne stilles til rådighed til forskningsbrug og udgør ikke professionel rådgivning.

4) Privatliv og moderering. For at bevare sikkerhed og tillid kan tider og koordinater blive gjort mindre præcise, og tekniske metadata kan blive fjernet eller anonymiseret. Vi kan skjule eller slette indhold, som efter vores rimelige vurdering er spam, forfalsket, ulovligt eller forstyrrende. Ærlige målinger behandler vi med omtanke og forsøger at bevare.

5) Cookies. Webstedet anvender kun en kortlivet teknisk sessionscookie; den forsvinder, når dit besøg slutter. Vi gemmer ingen andre spor.

Kære venner, kortet er frugten af fælles indsats og åbne hjerter. Se det som en skitse over terrænet, ikke en millimeternøjagtig plan. Hvis arbejdet taler til dig, så vær med – sammen kan vi gøre det endnu bedre.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'legal_title', 'Juridiske oplysninger') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'legend_attention', 'OBS') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'legend_button_tooltip', 'Åbn forklaringen over strålingsniveauer.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'legend_danger', 'Fare') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'legend_full_da', 'Denne skala viser, hvor sikker et sted er for liv, vand og mad.
Husk: målinger kan være ufuldstændige; brug den kun som vejledning.

Grøn (0–11 µR/h)
Næsten naturlig baggrund.
• Brøndvand er som regel sikkert.
• Planter kan dyrkes uden test.

Gul (11–30 µR/h)
Forhøjet baggrund.
• Tjek vand og jord.
• Undersøg mad før du spiser.

Rød (30–100 µR/h)
Alvorlig forurening.
• Drik ikke vandet.
• Dyrkning eller spisning her er risikabelt; laboratorietest kræves.

Sort (>100 µR/h)
Kritisk område.
• Vand og mad kan ikke bruges.
• Kun korte ophold med beskyttelse.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'legend_safe', 'Sikker') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'legend_title', 'Forklaring') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'license_full', 'Dette projekt drives under <a href="/LICENSE" target="_blank">MIT License</a>. Den fulde tekst ligger i roden af repositoriet og på MIT''s hjemmeside. Du må studere, dele og ændre koden, så længe disse friheder følger dit arbejde. Forskningsdata udgives under Creative Commons 1.0, så målingerne forbliver i det offentlige domæne.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'license_title', 'Licens') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_chart_all', 'Alle målinger') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_chart_averaged', 'Gennemsnit over [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_chart_close', 'Luk') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_chart_day', 'Sidste 24 timer') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_chart_month', 'Sidste 30 dage') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_chart_link', 'Åbn strålingsdiagrammer') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_country', 'Land') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_desc_generic', 'Denne Safecast-sensor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_desc_location', 'rapporterer fra [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_desc_location_unknown', 'rapporterer fra et ukendt område') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_desc_named', 'Enhed [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_desc_no_tube', 'overvåger strålingsniveauer.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_desc_radiation_sensor', 'Dette er en strålingssensor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_desc_transport_air', 'under flyvning') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_desc_transport_bike', 'på cykel') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_desc_transport_car', 'i bil') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_desc_transport_unknown', 'stationær') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_desc_transport_walk', 'til fods') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_desc_tube', 'med [[tube]]-detektor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_detector', 'Detektor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_device_id', 'Enheds-ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_device_label', 'Enhed') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_extra_intro', 'Miljø') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_history_error', 'Historik kan ikke indlæses lige nu.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_humidity', 'Luftfugtighed') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_last_seen', 'Seneste måling') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_local_time', 'Lokal tid nu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_no_data', 'Ingen data registreret i denne periode.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_pressure', 'Lufttryk') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_temperature', 'Lufttemperatur') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_temperature_f', 'Lufttemperatur (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_title', 'Safecast realtidssensor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_transport', 'Bevægelse') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_transport_air', 'Fly') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_transport_bike', 'Cykel eller løbehjul') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_transport_car', 'Bil eller varevogn') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_transport_unknown', 'Ikke registreret') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'live_marker_transport_walk', 'Til fods') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'locate_button_tooltip', 'Centrér kortet på min placering') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'location_error', 'Der opstod en fejl under hentning af placeringen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'location_permission_denied', 'Adgang til placering nægtet.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'location_timeout', 'Tidsfrist for placeringsanmodning udløb.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'location_unavailable', 'Placering ikke tilgængelig.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'no_results_found', 'Ingen resultater fundet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'processing_complete', 'Behandling fuldført!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'processing_on_server', 'Behandler på serveren...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'qr_button_tooltip', 'QR-kode til linket for dette kortområde.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'radiation_dose', 'Dosisrate') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'search_error', 'Søgefejl. Prøv venligst igen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'search_input_tooltip', 'Søg efter et sted ved at skrive de første bogstaver. En liste med forslag vises.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'select_files', 'Vælg mindst én fil') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'short_link_tooltip', 'Klik for at kopiere et kort delingslink') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'sources_full', 'Vi takker alle, der deler målinger.

Anonyme uploads tegner stille spor på kortet.
<a href="https://safecast.org" target="_blank">Safecast</a> plejer et globalt arkiv af aflæsninger.
<a href="https://atomfast.net" target="_blank">Atomfast</a> holder Atomcloud i gang.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> samler indsigter fra Radiacode.

Hvert bidrag udvider det fælles billede; du er hjerteligt velkommen til at tilføje dit eget.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'sources_title', 'Datakilder') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'speed', 'Hastighed') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'speed_filter_tooltip_accuracy', 'Langsommere målinger ligger tættest på jorden, så fodgængerdata er mest præcise.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'speed_filter_tooltip_car', 'Bil: 7–200 km/t til kørsler og mobile målinger.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'speed_filter_tooltip_intro', 'Vælg hvilke målinger der vises efter bevægelseshastighed.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'speed_filter_tooltip_live', 'Safecast-hjerte: data i realtid fra safecast.org. Brug afkrydsningsfeltet til at vise eller skjule live-aflæsninger.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'speed_filter_tooltip_ped', 'Fodgænger: under 7 km/t til gang eller stationære målinger.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'speed_filter_tooltip_plane', 'Fly: 200+ km/t til målinger fra luften.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'speed_filter_tooltip_title', 'Hastighedsfiltre') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'theme_toggle_tooltip', 'Skift mellem det lyse og mørke korttema.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'title', 'Safecasts isotopkort — Radiologisk kort over planeten Jorden') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'track_id', 'Spor-ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'upload_button', 'Upload [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'upload_button_tooltip', 'Tilføj din målerute til kortet. Understøttede formater: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Du kan uploade flere filer, og efter upload åbnes sporesiden.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'upload_error', 'Fejl') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'waiting_for_server', 'Venter på serveren...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('da', 'your_location', 'Din placering') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_archive_desc', 'Lädt ein tgz-Paket mit allen veröffentlichten .json-Dateien, sofern das JSON-Archiv aktiv ist.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_archive_link', 'Wöchentliches Archiv herunterladen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_archive_note', 'Ist das Archiv deaktiviert, liefert der Server HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_archive_title', 'Wöchentliches Archivpaket') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_latest_desc', 'Gibt die neuesten Messpunkte nahe der angegebenen Breite, Länge und Suchradius in Metern zurück.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_latest_link', 'Neueste bei Tokio') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_latest_note', 'Passen Sie lat, lon und radius_m an, um Ihre Region zu fokussieren.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_latest_title', 'Neueste Messungen in der Nähe') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_root_desc', 'Listet Metadaten, Datensatzstatistiken und Links zu allen weiteren Endpunkten auf.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_root_link', 'Öffnen Sie /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_root_note', 'Starten Sie hier, um Sammlungen und den Serverstatus zu entdecken.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_root_title', 'Maschinenlesbarer Index') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_track_index_desc', 'Lädt die N-te Strecke und liefert dasselbe JSON wie /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_track_index_link', 'Index 1 abrufen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_track_index_note', 'Ersetzen Sie die Zahl, um einen anderen Eintrag zu holen. Antworten streamen als zeilenweises JSON.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_track_index_title', 'Strecke nach Index herunterladen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_tracks_desc', 'Streamt einen Katalog veröffentlichter Strecken inklusive Namen, Anzahl und Download-URLs.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_tracks_link', 'Erste Strecken listen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_tracks_note', 'Nutzen Sie die Parameter limit und offset, um lange Listen zu paginieren.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_example_tracks_title', 'Alle Streckenzusammenfassungen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_examples_heading', 'Ausgewählte Endpunkte') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_examples_note', 'Alle Links liefern JSON. Beim Streamen können Browser große Dateien herunterladen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_intro', 'Die JSON-API spiegelt die Daten wider, die auf der Karte angezeigt werden. Jedes Beispiel öffnet sich in einem neuen Tab, damit Sie die Rohantwort prüfen können.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_more_docs', 'Benötigen Sie ausführlichere Informationen?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_more_docs_link_label', 'Ausführliche Anleitung öffnen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'api_title', 'API-Schnellstart') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'attr_legal', 'Rechtliche Hinweise') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'attr_license', 'Lizenz') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'attr_sources', 'Datenquellen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'back_to_all_tracks', 'Zurück zur kombinierten Track-Karte.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'count_rate', 'Zählrate') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'date_slider_tooltip', 'Messwerte nach Datum oder Jahr filtern.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'date_time', 'Datum und Uhrzeit') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'description', 'Safecasts Strahlungskarte') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'download_track_cim', 'Track herunterladen (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'duration_days', '[[count]] Tage') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'duration_hours', '[[count]] Std') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'duration_minutes', '[[count]] Min') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'duration_months', '[[count]] Monate') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'duration_weeks', '[[count]] Wochen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'error_during_upload', 'Fehler beim Hochladen!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'error_processing_files', 'Fehler bei der Verarbeitung der Dateien!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'files_uploaded', 'Dateien hochgeladen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'geolocation_not_supported', 'Ihr Browser unterstützt keine Geolokalisierung.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'github_link_tooltip_desc', 'Open-Source-Projekt aus der Community.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'github_link_tooltip_title', 'Safecast New Map auf GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'github_link_tooltip_version', 'Aktuelle Version: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'home_explore_global', 'Überspringen und Weltkarte erkunden') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'home_or', 'oder') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'home_search_placeholder', 'Stadt, Region oder Land eingeben...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'home_subtitle', 'Geben Sie einen Ort ein, um zu beginnen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'home_title', 'Safecast Strahlungskarte') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'home_use_location', 'Meinen Standort verwenden') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'legal_contact', 'Für Rückmeldungen melde dich unter:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'legal_full', 'Präambel. Wir bauen eine offene Karte, auf der Menschen aus aller Welt ihre Dosimetermessungen zum Wohl aller teilen – für Wissenschaft, Umwelt, Bildung und Sicherheit. Mit jeder Veröffentlichung hilfst du vielen. Bitte geh sorgsam mit diesem Gemeinschaftsprojekt um.

1) Verantwortung. Für Richtigkeit und Inhalt der Daten, die du einreichst, trägst du selbst die Verantwortung. Die Angaben werden „wie gesehen“ veröffentlicht und genutzt. Der Dienst bestätigt oder garantiert weder Korrektheit noch Vollständigkeit oder Eignung für irgendeinen Zweck und übernimmt keine Haftung für Folgen ihrer Verwendung.

2) Offenheit und Lizenz. Wenn du Messreihen, Zeitpunkte, ungefähre Orte, Gerätemodelle oder andere Details bereitstellst, weißt du, dass sie allen offenstehen und frei unter der Datenlizenz CC0 1.0 (Public Domain) genutzt werden dürfen. Der Programmcode bleibt unter der MIT-Lizenz verfügbar. Die Urheberschaft bleibt bei dir; es gibt keine Vergütung; eine Weitergabe durch Dritte entzieht sich unserem Einfluss.

3) „Wie gesehen“ und ohne Vorprüfung. Beiträge werden ohne vorherige Kontrolle veröffentlicht. Wir können weder eine kalibrierte Messtechnik noch Fehlerfreiheit garantieren. Die Informationen dienen der Forschung und stellen keine fachliche Empfehlung dar.

4) Datenschutz und Moderation. Um Sicherheit und Vertrauen zu wahren, können Zeitangaben und Koordinaten verallgemeinert werden; technische Metadaten dürfen entfernt oder anonymisiert werden. Inhalte, die nach unserem vernünftigen Ermessen Spam, Fälschung, rechtswidrig oder störend sind, dürfen wir ausblenden oder löschen. Mit aufrichtig erhobenen Messungen gehen wir behutsam um und versuchen, sie zu bewahren.

5) Cookies. Diese Website verwendet lediglich ein kurzlebiges technisches Sitzungs-Cookie; es erlischt mit dem Ende deines Besuchs. Weitere Spuren speichern wir nicht.

Freundinnen und Freunde, diese Karte ist das Ergebnis gemeinsamer Arbeit und offener Herzen. Betrachtet sie als Skizze der Landschaft, nicht als exakten Plan. Wenn euch unser Engagement bewegt, macht mit – gemeinsam machen wir sie besser.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'legal_title', 'Rechtliches') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'legend_attention', 'Achtung') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'legend_button_tooltip', 'Strahlungslegende öffnen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'legend_danger', 'Gefahr') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'legend_full_de', 'Diese Skala zeigt, wie sicher ein Ort für Menschen, Wasser und Nahrung ist.
Denk daran: Messungen können unvollständig sein; nutze sie nur als Richtwert.

Grün (0–11 µR/h)
Natürliche Hintergrundstrahlung.
• Brunnenwasser meist sicher.
• Pflanzen können ohne Tests wachsen.

Gelb (11–30 µR/h)
Erhöhte Werte.
• Wasser und Boden prüfen.
• Jede Ernte vor dem Essen testen.

Rot (30–100 µR/h)
Starke Belastung.
• Wasser nicht trinken.
• Anbau oder Verzehr hier gefährlich; Laborprüfungen nötig.

Schwarz (>100 µR/h)
Kritische Zone.
• Wasser und Nahrung unbrauchbar.
• Längerer Aufenthalt verboten; nur kurz mit Schutz.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'legend_safe', 'Sicher') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'legend_title', 'Legende') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'license_full', 'Dieses Projekt steht unter der <a href="/LICENSE" target="_blank">MIT License</a>. Der vollständige Text liegt im Wurzelverzeichnis des Repos und auf der MIT-Seite. Du darfst den Code studieren, weitergeben und verändern, solange diese Freiheiten mit deiner Arbeit weitergegeben werden. Die Forschungsdaten erscheinen unter der Lizenz Creative Commons 1.0, damit die Messwerte gemeinfrei bleiben.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'license_title', 'Lizenz') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_chart_all', 'Alle Messwerte') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_chart_averaged', 'Durchschnitt über [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_chart_close', 'Schließen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_chart_day', 'Letzte 24 Stunden') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_chart_month', 'Letzte 30 Tage') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_chart_link', 'Strahlungsdiagramme öffnen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_country', 'Land') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_desc_generic', 'Dieser Safecast-Sensor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_desc_location', 'meldet aus [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_desc_location_unknown', 'meldet aus einem unbekannten Gebiet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_desc_named', 'Gerät [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_desc_no_tube', 'überwacht die Strahlungswerte.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_desc_radiation_sensor', 'Dies ist ein Strahlungssensor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_desc_transport_air', 'im Flug') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_desc_transport_bike', 'per Fahrrad') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_desc_transport_car', 'per Auto') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_desc_transport_unknown', 'stationär') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_desc_transport_walk', 'zu Fuß') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_desc_tube', 'mit [[tube]]-Detektor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_detector', 'Detektor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_device_id', 'Geräte-ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_device_label', 'Gerät') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_extra_intro', 'Umgebung') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_history_error', 'Verlauf kann derzeit nicht geladen werden.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_humidity', 'Luftfeuchtigkeit') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_last_seen', 'Letzte Messung') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_local_time', 'Ortszeit jetzt') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_no_data', 'Keine Daten in diesem Zeitraum erfasst.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_pressure', 'Luftdruck') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_temperature', 'Lufttemperatur') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_temperature_f', 'Lufttemperatur (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_title', 'Safecast Echtzeit-Sensor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_transport', 'Bewegung') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_transport_air', 'Flugzeug') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_transport_bike', 'Fahrrad oder Roller') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_transport_car', 'Auto oder Transporter') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_transport_unknown', 'Nicht erkannt') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'live_marker_transport_walk', 'Zu Fuß') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'locate_button_tooltip', 'Karte auf meinen Standort zentrieren') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'location_error', 'Fehler beim Abrufen des Standorts.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'location_permission_denied', 'Zugriff auf den Standort verweigert.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'location_timeout', 'Standortanfrage abgelaufen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'location_unavailable', 'Standort nicht verfügbar.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'no_results_found', 'Keine Ergebnisse gefunden') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'processing_complete', 'Verarbeitung abgeschlossen!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'processing_on_server', 'Verarbeitung auf dem Server...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'qr_button_tooltip', 'QR-Code des Links zu diesem Kartenausschnitt.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'radiation_dose', 'Dosisleistung') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'search_error', 'Suchfehler. Bitte versuchen Sie es erneut.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'search_input_tooltip', 'Suchen Sie nach einem Ort, indem Sie die ersten Buchstaben eingeben. Eine Vorschlagsliste wird angezeigt.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'select_files', 'Bitte wählen Sie mindestens eine Datei aus') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'short_link_tooltip', 'Klicken, um einen kurzen Freigabelink zu kopieren') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'sources_full', 'Wir danken allen, die Messungen teilen.

Anonyme Uploads zeichnen leise Wege auf die Karte.
<a href="https://safecast.org" target="_blank">Safecast</a> pflegt ein weltweites Archiv von Messwerten.
<a href="https://atomfast.net" target="_blank">Atomfast</a> hält die Atomcloud am Laufen.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> sammelt Erkenntnisse aus Radiacode.

Jeder Beitrag erweitert das gemeinsame Bild; wir laden dich herzlich ein, selbst mitzuwirken.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'sources_title', 'Datenquellen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'speed', 'Geschwindigkeit') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'speed_filter_tooltip_accuracy', 'Langsamere Messungen bleiben näher am Boden, daher sind Fußgängerdaten am genauesten.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'speed_filter_tooltip_car', 'Auto: 7–200 km/h für Fahrten und mobile Messungen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'speed_filter_tooltip_intro', 'Wählen Sie, welche Messungen nach Reisegeschwindigkeit angezeigt werden.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'speed_filter_tooltip_live', 'Safecast-Herz: Echtzeitdaten von safecast.org. Nutzen Sie die Option, um Live-Messungen ein- oder auszublenden.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'speed_filter_tooltip_ped', 'Zu Fuß: unter 7 km/h für Geh- oder stationäre Messungen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'speed_filter_tooltip_plane', 'Flugzeug: 200+ km/h für Luftmessungen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'speed_filter_tooltip_title', 'Geschwindigkeitsfilter') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'theme_toggle_tooltip', 'Zwischen hellem und dunklem Kartenthema wechseln.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'title', 'Safecasts Isotopenkarte — Radiologische Karte des Planeten Erde') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'track_id', 'Track-ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'upload_button', 'Hochladen [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'upload_button_tooltip', 'Fügen Sie Ihre Messstrecke zur Karte hinzu. Unterstützte Formate: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Sie können mehrere Dateien hochladen, danach öffnet sich die Track-Seite.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'upload_error', 'Fehler') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'waiting_for_server', 'Warten auf den Server...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('de', 'your_location', 'Ihr Standort') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_archive_desc', 'Κατεβάζει ένα αρχείο tgz με όλα τα δημοσιευμένα αρχεία .json όταν είναι ενεργό το JSON αρχείο.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_archive_link', 'Λήψη εβδομαδιαίου αρχείου') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_archive_note', 'Αν το αρχείο είναι απενεργοποιημένο, ο διακομιστής επιστρέφει HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_archive_title', 'Εβδομαδιαίο πακέτο αρχείου') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_latest_desc', 'Επιστρέφει τα πιο πρόσφατα σημεία κοντά στα δοσμένα γεωγραφικά πλάτη, μήκη και την ακτίνα σε μέτρα.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_latest_link', 'Πρόσφατες γύρω από το Τόκιο') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_latest_note', 'Ρυθμίστε τα lat, lon και radius_m για τη δική σας περιοχή.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_latest_title', 'Πιο πρόσφατες μετρήσεις κοντά') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_root_desc', 'Παραθέτει μεταδεδομένα, στατιστικά συνόλου δεδομένων και συνδέσμους για κάθε άλλο τελικό σημείο.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_root_link', 'Άνοιγμα /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_root_note', 'Ξεκινήστε από εδώ για να ανακαλύψετε συλλογές και την κατάσταση του διακομιστή.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_root_title', 'Ευρετήριο αναγνώσιμο από μηχανή') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_track_index_desc', 'Επιστρέφει την N-οστή διαδρομή και δίνει το ίδιο JSON με το /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_track_index_link', 'Λήψη δείκτη 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_track_index_note', 'Αντικαταστήστε τον αριθμό για να φέρετε άλλη εγγραφή. Οι αποκρίσεις μεταδίδονται ως JSON ανά γραμμή.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_track_index_title', 'Λήψη διαδρομής κατά δείκτη') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_tracks_desc', 'Μεταδίδει έναν κατάλογο δημοσιευμένων διαδρομών με ονόματα, πλήθος και συνδέσμους λήψης.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_tracks_link', 'Λίστα πρώτων διαδρομών') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_tracks_note', 'Χρησιμοποιήστε τις παραμέτρους limit και offset για να περιηγηθείτε σε μεγάλες λίστες.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_example_tracks_title', 'Σύνοψη όλων των διαδρομών') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_examples_heading', 'Προβεβλημένα τελικά σημεία') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_examples_note', 'Όλοι οι σύνδεσμοι απαντούν με JSON. Όταν τα τελικά σημεία μεταδίδουν δεδομένα, ο περιηγητής μπορεί να κατεβάσει μεγάλα αρχεία.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_intro', 'Το JSON API αντικατοπτρίζει τα δεδομένα που εμφανίζονται στον χάρτη. Κάθε παράδειγμα ανοίγει σε νέα καρτέλα για να ελέγξετε την ακατέργαστη απόκριση.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_more_docs', 'Χρειάζεστε αναλυτικότερη τεκμηρίωση;') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_more_docs_link_label', 'Άνοιγμα αναλυτικού οδηγού') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'api_title', 'Γρήγορη έναρξη API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'attr_legal', 'Νομικές πληροφορίες') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'attr_license', 'Άδεια') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'attr_sources', 'Πηγές δεδομένων') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'back_to_all_tracks', 'Επιστροφή στον συνδυαστικό χάρτη διαδρομών.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'count_rate', 'Ρυθμός καταμέτρησης') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'date_slider_tooltip', 'Φιλτράρετε μετρήσεις ανά ημερομηνία ή εύρος ετών.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'date_time', 'Ημερομηνία και ώρα') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'description', 'Χάρτης ακτινοβολίας της Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'download_track_cim', 'Λήψη διαδρομής (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'duration_days', '[[count]] ημέρες') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'duration_hours', '[[count]] ώρ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'duration_minutes', '[[count]] λεπ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'duration_months', '[[count]] μήνες') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'duration_weeks', '[[count]] εβδομάδες') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'error_during_upload', 'Σφάλμα κατά τη μεταφόρτωση!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'error_processing_files', 'Σφάλμα κατά την επεξεργασία των αρχείων!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'files_uploaded', 'Τα αρχεία μεταφορτώθηκαν') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'geolocation_not_supported', 'Το πρόγραμμα περιήγησής σας δεν υποστηρίζει γεωεντοπισμό.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'github_link_tooltip_desc', 'Έργο ανοικτού κώδικα που συντηρείται από την κοινότητα.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'github_link_tooltip_title', 'Safecast New Map στο GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'github_link_tooltip_version', 'Τρέχουσα έκδοση: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'home_explore_global', 'Παράλειψη και εξερεύνηση του παγκόσμιου χάρτη') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'home_or', 'ή') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'home_search_placeholder', 'Εισάγετε πόλη, περιοχή ή χώρα...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'home_subtitle', 'Εισάγετε μια τοποθεσία για να ξεκινήσετε.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'home_title', 'Ραδιολογικός χάρτης Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'home_use_location', 'Χρήση της τοποθεσίας μου') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'legal_contact', 'Για σχόλια επικοινωνήστε στο:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'legal_full', 'Προοίμιο. Δημιουργούμε έναν ανοιχτό χάρτη όπου άνθρωποι απ’ όλον τον κόσμο μοιράζονται μετρήσεις δοσίμετρων για το κοινό καλό — την επιστήμη, το περιβάλλον, την εκπαίδευση και την ασφάλεια. Δημοσιεύοντας τα δεδομένα σας βοηθάτε πολλούς. Σας ζητούμε να φροντίζετε αυτό το συλλογικό έργο.

1) Ευθύνη. Η ακρίβεια και το περιεχόμενο όσων υποβάλλετε παραμένουν δική σας ευθύνη. Τα δεδομένα δημοσιεύονται και χρησιμοποιούνται «ως έχουν». Η υπηρεσία δεν επιβεβαιώνει ούτε εγγυάται την ορθότητά τους, την πληρότητα ή την καταλληλότητά τους για οποιονδήποτε σκοπό και δεν ευθύνεται για τυχόν συνέπειες από τη χρήση τους.

2) Ανοιχτότητα και άδεια. Μοιραζόμενοι μετρήσεις, ημερομηνίες, προσεγγιστικές τοποθεσίες, μοντέλα συσκευών ή άλλα στοιχεία κατανοείτε ότι θα γίνουν διαθέσιμα σε όλους και μπορούν να χρησιμοποιηθούν ελεύθερα με βάση την άδεια CC0 1.0 (Public Domain) για τα δεδομένα. Ο κώδικας παραμένει υπό την άδεια MIT. Η πατρότητα παραμένει δική σας· δεν προβλέπεται αμοιβή· η περαιτέρω αναδιανομή από τρίτους βρίσκεται εκτός του ελέγχου μας.

3) «Ως έχει» και χωρίς επαλήθευση. Οι δημοσιεύσεις εμφανίζονται χωρίς προληπτικό έλεγχο. Δεν μπορούμε να εγγυηθούμε τη βαθμονόμηση των οργάνων ούτε την απουσία σφαλμάτων. Οι πληροφορίες παρέχονται για ερευνητικούς σκοπούς και δεν αποτελούν επαγγελματική σύσταση.

4) Ιδιωτικότητα και επιμέλεια. Για να διατηρήσουμε ασφάλεια και εμπιστοσύνη, μπορεί να γενικεύουμε χρόνους και συντεταγμένες, ενώ τεχνικά μεταδεδομένα ενδέχεται να αφαιρούνται ή να ανωνυμοποιούνται. Διατηρούμε το δικαίωμα να αποκρύπτουμε ή να διαγράφουμε υλικό που, κατά την εύλογη κρίση μας, αποτελεί ανεπιθύμητη αλληλογραφία, πλαστογραφία, παραβίαση νόμου ή ενόχληση της υπηρεσίας. Τιμούμε τις καλόπιστες μετρήσεις και προσπαθούμε να τις προστατεύουμε.

5) Cookies. Ο ιστότοπος χρησιμοποιεί μόνο ένα τεχνικό cookie συνεδρίας μικρής διάρκειας· εξαφανίζεται με τη λήξη της επίσκεψής σας. Δεν διατηρούμε άλλα ίχνη.

Φίλες και φίλοι, αυτός ο χάρτης είναι καρπός συλλογικής προσπάθειας και ανοιχτών καρδιών. Δείτε τον ως ένα πρόχειρο περίγραμμα του τόπου και όχι ως ακριβές σχέδιο. Αν το έργο μας σας αγγίζει, ελάτε μαζί μας — μαζί μπορούμε να το εξελίξουμε.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'legal_title', 'Νομικές πληροφορίες') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'legend_attention', 'Προσοχή') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'legend_button_tooltip', 'Άνοιγμα του υπομνήματος επιπέδων ακτινοβολίας.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'legend_danger', 'Κίνδυνος') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'legend_full_el', 'Αυτή η κλίμακα δείχνει πόσο ασφαλής είναι μια περιοχή για ζωή, νερό και τροφή.
Θυμήσου: οι μετρήσεις μπορεί να μην είναι πλήρεις· χρησιμοποίησέ τες μόνο ως οδηγό.

Πράσινο (0–11 µR/h)
Σχεδόν φυσικό υπόβαθρο.
• Το νερό πηγών είναι γενικά ασφαλές.
• Μπορείς να καλλιεργείς χωρίς έλεγχο.

Κίτρινο (11–30 µR/h)
Αυξημένο υπόβαθρο.
• Έλεγξε νερό και έδαφος.
• Εξέτασε κάθε τρόφιμο πριν το φας.

Κόκκινο (30–100 µR/h)
Σοβαρή ρύπανση.
• Μην πίνεις το νερό.
• Η καλλιέργεια ή κατανάλωση εδώ είναι επικίνδυνη· απαιτούνται εργαστηριακοί έλεγχοι.

Μαύρο (>100 µR/h)
Κρίσιμη ζώνη.
• Νερό και τροφή ακατάλληλα.
• Μόνο σύντομη παραμονή με προστασία.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'legend_safe', 'Ασφαλές') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'legend_title', 'Υπόμνημα') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'license_full', 'Το έργο αυτό εξελίσσεται υπό την άδεια <a href="/LICENSE" target="_blank">MIT License</a>. Το πλήρες κείμενο βρίσκεται στη ρίζα του αποθετηρίου και στον ιστότοπο του MIT. Μπορείτε να μελετήσετε, να μοιραστείτε και να τροποποιήσετε τον κώδικα, αρκεί αυτές οι ελευθερίες να συνοδεύουν και τη δική σας δουλειά. Τα ερευνητικά δεδομένα δημοσιεύονται υπό την άδεια Creative Commons 1.0 ώστε οι μετρήσεις να παραμένουν στο δημόσιο τομέα.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'license_title', 'Άδεια') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_chart_all', 'Όλες οι μετρήσεις') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_chart_averaged', 'Μέσος όρος σε [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_chart_close', 'Κλείσιμο') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_chart_day', 'Τελευταίες 24 ώρες') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_chart_month', 'Τελευταίες 30 ημέρες') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_chart_link', 'Άνοιγμα γραφημάτων ακτινοβολίας') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_country', 'Χώρα') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_desc_generic', 'Αυτός ο αισθητήρας Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_desc_location', 'αναφέρει από [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_desc_location_unknown', 'αναφέρει από άγνωστη περιοχή') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_desc_named', 'Συσκευή [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_desc_no_tube', 'παρακολουθεί τα επίπεδα ακτινοβολίας.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_desc_radiation_sensor', 'Αυτός είναι αισθητήρας ακτινοβολίας.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_desc_transport_air', 'εν πτήσει') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_desc_transport_bike', 'με ποδήλατο') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_desc_transport_car', 'με αυτοκίνητο') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_desc_transport_unknown', 'σταθερός') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_desc_transport_walk', 'με τα πόδια') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_desc_tube', 'με ανιχνευτή [[tube]].') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_detector', 'Ανιχνευτής') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_device_id', 'ID συσκευής') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_device_label', 'Συσκευή') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_extra_intro', 'Περιβάλλον') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_history_error', 'Αδυναμία φόρτωσης ιστορικού αυτή τη στιγμή.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_humidity', 'Υγρασία') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_last_seen', 'Τελευταία μέτρηση') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_local_time', 'Τοπική ώρα τώρα') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_no_data', 'Δεν καταγράφηκαν δεδομένα σε αυτή την περίοδο.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_pressure', 'Ατμοσφαιρική πίεση') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_temperature', 'Θερμοκρασία αέρα') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_temperature_f', 'Θερμοκρασία αέρα (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_title', 'Αισθητήρας Safecast σε πραγματικό χρόνο') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_transport', 'Κίνηση') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_transport_air', 'Αεροσκάφος') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_transport_bike', 'Ποδήλατο ή πατίνι') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_transport_car', 'Αυτοκίνητο ή βαν') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_transport_unknown', 'Δεν ανιχνεύτηκε') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'live_marker_transport_walk', 'Με τα πόδια') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'locate_button_tooltip', 'Κεντράρετε τον χάρτη στη θέση μου') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'location_error', 'Παρουσιάστηκε σφάλμα κατά την απόκτηση της τοποθεσίας.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'location_permission_denied', 'Δεν επιτρέπεται η πρόσβαση στην τοποθεσία.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'location_timeout', 'Έληξε το αίτημα τοποθεσίας.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'location_unavailable', 'Η τοποθεσία δεν είναι διαθέσιμη.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'no_results_found', 'Δεν βρέθηκαν αποτελέσματα') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'processing_complete', 'Η επεξεργασία ολοκληρώθηκε!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'processing_on_server', 'Επεξεργασία στον διακομιστή...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'qr_button_tooltip', 'QR κωδικός του συνδέσμου για αυτή την περιοχή του χάρτη.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'radiation_dose', 'Ρυθμός δόσης') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'search_error', 'Σφάλμα αναζήτησης. Δοκιμάστε ξανά.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'search_input_tooltip', 'Αναζητήστε ένα μέρος πληκτρολογώντας τα πρώτα γράμματα του ονόματός του.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'select_files', 'Παρακαλώ επιλέξτε τουλάχιστον ένα αρχείο') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'short_link_tooltip', 'Κάντε κλικ για αντιγραφή σύντομου συνδέσμου κοινής χρήσης') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'sources_full', 'Ευχαριστούμε όλους όσοι μοιράζονται μετρήσεις.

Ανώνυμες αποστολές χαράσσουν ήσυχες διαδρομές στον χάρτη.
<a href="https://safecast.org" target="_blank">Safecast</a> καλλιεργεί ένα παγκόσμιο αρχείο καταγραφών.
<a href="https://atomfast.net" target="_blank">Atomfast</a> κρατά το Atomcloud σε λειτουργία.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> συγκεντρώνει γνώσεις από τα Radiacode.

Κάθε συνεισφορά μεγαλώνει τη συλλογική εικόνα· σας προσκαλούμε θερμά να προσθέσετε και τη δική σας.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'sources_title', 'Πηγές δεδομένων') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'speed', 'Ταχύτητα') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'speed_filter_tooltip_accuracy', 'Οι πιο αργές μετρήσεις παραμένουν πιο κοντά στο έδαφος, οπότε τα δεδομένα πεζών είναι τα πιο ακριβή.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'speed_filter_tooltip_car', 'Αυτοκίνητο: 7–200 χλμ./ώρα για διαδρομές και κινητές μετρήσεις.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'speed_filter_tooltip_intro', 'Επιλέξτε ποιες μετρήσεις εμφανίζονται ανάλογα με την ταχύτητα κίνησης.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'speed_filter_tooltip_live', 'Καρδιά Safecast: δεδομένα σε πραγματικό χρόνο από το safecast.org. Χρησιμοποιήστε την επιλογή για να εμφανίσετε ή να κρύψετε τις ζωντανές μετρήσεις.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'speed_filter_tooltip_ped', 'Πεζός: κάτω από 7 χλμ./ώρα για μετρήσεις σε περπάτημα ή στάση.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'speed_filter_tooltip_plane', 'Αεροπλάνο: 200+ χλμ./ώρα για εναέριες καταγραφές.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'speed_filter_tooltip_title', 'Φίλτρα ταχύτητας') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'theme_toggle_tooltip', 'Εναλλάξτε μεταξύ φωτεινού και σκοτεινού θέματος χάρτη.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'title', 'Χάρτης ισοτόπων της Safecast — Ραδιολογικός χάρτης του πλανήτη Γη') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'track_id', 'Αναγνωριστικό διαδρομής') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'upload_button', 'Μεταφόρτωση [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'upload_button_tooltip', 'Προσθέστε τη διαδρομή μετρήσεών σας στον χάρτη. Υποστηριζόμενες μορφές: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Μπορείτε να μεταφορτώσετε πολλά αρχεία· μετά τη μεταφόρτωση θα ανοίξει η σελίδα της διαδρομής.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'upload_error', 'Σφάλμα') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'waiting_for_server', 'Αναμονή για τον διακομιστή...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('el', 'your_location', 'Η τοποθεσία σας') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_archive_desc', 'Downloads a tgz bundle with every published .json file when the server enables the JSON archive.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_archive_link', 'Download weekly archive') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_archive_note', 'If archives are disabled the server returns HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_archive_title', 'Weekly archive bundle') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_latest_desc', 'Returns the newest markers near the provided latitude, longitude, and search radius in meters.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_latest_link', 'Latest near Tokyo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_latest_note', 'Change the coordinates and radius_m to focus on your own area.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_latest_title', 'Latest nearby measurements') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_root_desc', 'Lists metadata, dataset statistics, and hyperlinks to every other endpoint.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_root_link', 'Open /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_root_note', 'Start here to discover collections and server status.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_root_title', 'Machine-readable index') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_track_index_desc', 'Resolves the Nth track and returns the same JSON document as /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_track_index_link', 'Fetch index 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_track_index_note', 'Replace the number to fetch a different entry. Responses stream as newline-delimited JSON.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_track_index_title', 'Download track by index') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_tracks_desc', 'Streams a catalog of published tracks including names, counts, and download URLs.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_tracks_link', 'List first tracks') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_tracks_note', 'Use limit and offset parameters to paginate through long lists.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_example_tracks_title', 'All track summaries') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_examples_heading', 'Featured endpoints') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_examples_note', 'All links respond with JSON. Browsers may download large files when endpoints stream data.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_intro', 'The JSON API mirrors the data shown on the map. Each example opens in a new tab so you can inspect the raw response.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_more_docs', 'Need deeper coverage?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_more_docs_link_label', 'Open detailed guide') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'api_title', 'API quickstart') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'attr_legal', 'Legal') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'attr_license', 'License') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'attr_sources', 'Data Sources') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'back_to_all_tracks', 'Back to the combined track map.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'count_rate', 'Count Rate') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'date_slider_tooltip', 'Filter measurements by date or year range.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'date_time', 'Date and Time') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'description', 'Safecast New Map') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'download_track_cim', 'Download track (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'duration_days', '[[count]] days') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'duration_hours', '[[count]] hr') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'duration_minutes', '[[count]] min') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'duration_months', '[[count]] months') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'duration_weeks', '[[count]] weeks') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'error_during_upload', 'Error during upload!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'error_processing_files', 'Error processing files!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'files_uploaded', 'Files uploaded') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'geolocation_not_supported', 'Your browser does not support geolocation.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'github_link_tooltip_desc', 'Open-source project maintained by the community.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'github_link_tooltip_title', 'Safecast New Map on GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'github_link_tooltip_version', 'Current version: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'home_explore_global', 'Skip and explore the global map') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'home_or', 'or') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'home_search_placeholder', 'Enter a city, region, or country...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'home_subtitle', 'Enter a location to begin.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'home_title', 'Safecast Radiological Map') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'home_use_location', 'Use my location') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'legal_contact', 'For feedback please contact:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'legal_full', 'Preamble. We are building an open map where people around the globe share dosimeter readings for the common good — science, ecology, education, and safety. By publishing your data you help many. Please treat this shared work with care.

1) Responsibility. The accuracy and content of any information you submit remain your responsibility. Data are published and used “as is”. The service does not confirm or guarantee their correctness, completeness, or suitability for any purpose and bears no responsibility for any possible outcome of their use.

2) Openness and license. By sharing dosimetry readings, dates, approximate locations, device models, or other facts, you understand they become available to everyone and may be used freely under the CC0 1.0 (Public Domain) license for data, while the code remains available under the MIT license. Authorship stays with you; no compensation is provided; further redistribution by third parties is beyond our control.

3) “As is” and no verification. Publications appear without prior review. We cannot guarantee instrument calibration or the absence of mistakes. The information is shared for research and is not a professional recommendation.

4) Privacy and moderation. To keep the map safe and trustworthy, timestamps and coordinates may be generalized, and technical metadata may be removed or anonymized. We may hide or delete materials that, in our reasonable opinion, are spam, falsified, unlawful, or disruptive. We handle bona fide measurements with care and strive to preserve them.

5) Cookies. The site relies on a short-lived technical session cookie; it disappears when your visit ends. We do not keep any other traces.

Friends, this map is the result of shared effort and open hearts. Please treat it as a sketch of the terrain, not a precise blueprint. If our work resonates with you, join in — together we can make it better.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'legal_title', 'Legal') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'legend_attention', 'Attention') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'legend_button_tooltip', 'Open the radiation safety legend.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'legend_danger', 'Danger') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'legend_full_en', '<h1>📍 How to read a radiation map</h1><section style="margin:0 0 12px;"><h2>The key points from the outset</h2><ul style="margin:6px 0 0 18px;"><li>The number shown on the map reflects <strong>only part of the radiation hazard</strong> — the portion the instrument can detect.</li><li>Real harm is more often caused <strong>not by sudden spikes</strong> but by <strong>quiet, long-term exposure</strong> that builds up over years.</li><li>Even in a "green" area, risk may still exist due to <strong>invisible sources</strong>: radon, water, food and air.</li><li>A map is <strong>a guide, not a guarantee</strong>. It helps you think but <strong>does not replace common sense</strong>.</li></ul></section><hr /><section style="margin:0 0 12px; display:flex; gap:10px; align-items:flex-start;"><div style="min-width:84px; padding:8px 10px; border-radius:10px; background:#1b5e20; color:#fff; text-align:center; font-weight:700;">🟩<br />0–11</div><div><h2 style="margin:0;">Green zone</h2><h3 style="margin:4px 0; font-weight:600;">"Generally safe — but stay attentive"</h3><p style="margin:4px 0;"><strong>0–11 microroentgens per hour</strong> — close to natural background.</p><ul style="margin:4px 0 0 18px;"><li>you can usually live here permanently</li><li>drink the water</li><li>grow and eat local produce</li><li>this <strong>does not mean radiation is absent</strong></li><li>it means levels are <strong>about what the body is used to</strong></li><li>⚠ hidden risks still exist: indoor radon; radioactive impurities in water; accumulation of radionuclides in food and tobacco</li></ul></div></section><hr /><section style="margin:0 0 12px; display:flex; gap:10px; align-items:flex-start;"><div style="min-width:84px; padding:8px 10px; border-radius:10px; background:#f9a825; color:#111; text-align:center; font-weight:700;">🟨<br />11–30</div><div><h2 style="margin:0;">Yellow zone</h2><h3 style="margin:4px 0; font-weight:600;">"Pay attention and find out why"</h3><p style="margin:4px 0;"><strong>11–30 microroentgens per hour</strong>. Usually no immediate danger, but a yellow patch inside green is <strong>a reason to pause</strong>.</p><ul style="margin:4px 0 0 18px;"><li>possible causes: radon; alpha or beta contamination the dosimeter only partially sees; localised sources</li><li>sensible steps: test water and food; identify the source; measure indoor radon</li><li>harm <strong>can build up over time</strong>, even if you feel fine now</li></ul></div></section><hr /><section style="margin:0 0 12px; display:flex; gap:10px; align-items:flex-start;"><div style="min-width:84px; padding:8px 10px; border-radius:10px; background:#c62828; color:#fff; text-align:center; font-weight:700;">🟥<br />30–100</div><div><h2 style="margin:0;">Red zone</h2><h3 style="margin:4px 0; font-weight:600;">"Without checks, it is not worth the risk"</h3><p style="margin:4px 0;"><strong>30–100 microroentgens per hour</strong> — no longer a normal level.</p><ul style="margin:4px 0 0 18px;"><li>avoid: living here permanently; regularly drinking local water; relying on local food</li><li>possible causes: man-made radionuclides; contaminated soil; high radon in homes and water</li><li>recommended: drink water <strong>only after testing</strong>; eat food <strong>after checks</strong>; live here <strong>only with full awareness of the risks</strong></li></ul></div></section><hr /><section style="margin:0 0 12px; display:flex; gap:10px; align-items:flex-start;"><div style="min-width:84px; padding:8px 10px; border-radius:10px; background:#111; color:#fff; text-align:center; font-weight:700;">⬛<br />100+</div><div><h2 style="margin:0;">Black zone</h2><h3 style="margin:4px 0; font-weight:600;">"Extended exposure is unsafe"</h3><p style="margin:4px 0;"><strong>Over 100 microroentgens per hour</strong> — a serious level of exposure.</p><ul style="margin:4px 0 0 18px;"><li>risk fronts: external irradiation; radionuclides entering the body; potential skin damage</li><li>food and water may contain hazardous substances even if they look and smell normal</li><li>🔹 exceptions (brief skin exposure): air travel; high-altitude locations</li></ul></div></section><hr /><section style="margin:0 0 12px;"><h2>🔍 What a dosimeter shows — and what it misses</h2><div style="display:flex; gap:12px;"><div style="flex:1;"><h3 style="margin:4px 0;">Typically detected:</h3><ul style="margin:4px 0 0 18px;"><li>gamma radiation</li><li>sometimes part of beta radiation</li></ul></div><div style="flex:1;"><h3 style="margin:4px 0;">Often missed:</h3><ul style="margin:4px 0 0 18px;"><li>alpha radiation (especially dangerous once inside the body)</li><li>radon — a gas we breathe and drink with water</li><li>tritium in water and air</li><li>radionuclides in food</li><li>neutrons (rare, but extremely dangerous)</li><li>👉 a "green" reading <strong>does not equal full safety</strong>.</li></ul></div></div></section><hr /><section style="margin:0 0 12px;"><h2>🧬 How radiation actually enters the body</h2><div style="display:flex; gap:12px;"><div style="flex:1;"><h3 style="margin:4px 0;">The most dangerous routes are quiet and long</h3><ul style="margin:4px 0 0 18px;"><li><strong>Through breathing</strong> — radon, dust, smoke → lungs</li><li><strong>Through water</strong> — tritium and dissolved isotopes → whole body</li><li><strong>Through food</strong> — mushrooms, greens, milk → accumulation in tissues</li><li><strong>Through smoking</strong> — tobacco readily absorbs radionuclides</li></ul></div><div style="flex:1;"><h3 style="margin:4px 0;">Less dangerous:</h3><ul style="margin:4px 0 0 18px;"><li>brief external exposure</li><li>one-off contact</li></ul></div></div></section><hr /><section style="margin:0 0 12px;"><h2>⏳ Time is the main risk factor</h2><p style="margin:4px 0;">Radiation is dangerous <strong>because of prolonged exposure</strong>. Moderate doses taken continuously can be more harmful than a short spike. The body can eliminate radionuclides, but with constant exposure and other stresses those systems <strong>get overloaded</strong>.</p></section><hr /><section style="margin:0 0 12px;"><h2>👥 Who is affected more severely</h2><p style="margin:4px 0;">Higher risk for: older people; those with chronic illnesses. What one body tolerates without consequences may cause serious problems for another.</p></section><hr /><section style="margin:0 0 12px;"><h2>⚛️ Radioactive "disguise"</h2><p style="margin:4px 0;">Some radionuclides resemble ordinary substances to the body: <strong>tritium</strong> — like water; <strong>carbon-14</strong> — like normal carbon; <strong>potassium-40</strong> — like potassium; <strong>strontium, radium</strong> — like calcium. Cells <strong>do not recognise the substitution</strong> and incorporate them. Damage and mutations are possible even when the dosimeter reads "normal."</p></section><hr /><section style="margin:0;"><h2>🗺️ Why a map should not be trusted blindly</h2><p style="margin:4px 0;">Even a perfectly green map reflects <strong>only part of reality</strong>. Hidden factors are not shown. 👉 A map is a hint, not a pardon. You still need to think, check, and ask questions yourself.</p></section>') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'legend_safe', 'Safe') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'legend_title', 'Legend') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'license_full', 'This project grows under the <a href="/LICENSE" target="_blank">MIT License</a>. The full text rests in the repository''s root and on the MIT site. You may study, share, and alter the code, provided these freedoms travel with your work. Research datasets are released under Creative Commons 1.0 (CC0) so the measurements remain in the public domain.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'license_title', 'License') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_chart_all', 'All readings') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_chart_averaged', 'Averaged over [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_chart_close', 'Close') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_chart_day', 'Last 24 hours') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_chart_month', 'Last 30 days') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_chart_link', 'Open radiation charts') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_country', 'Country hint') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_desc_generic', 'This Safecast sensor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_desc_location', 'reports from [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_desc_location_unknown', 'reports from an unknown area') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_desc_named', 'Device [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_desc_no_tube', 'keeping watch on radiation levels.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_desc_radiation_sensor', 'This is a radiation sensor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_desc_transport_air', 'while flying') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_desc_transport_bike', 'while travelling by bike') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_desc_transport_car', 'while travelling by car') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_desc_transport_unknown', 'while staying in place') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_desc_transport_walk', 'while travelling on foot') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_desc_tube', 'using the [[tube]] detector.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_detector', 'Detector') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_device_id', 'Device ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_device_label', 'Device') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_extra_intro', 'Environment') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_history_error', 'Unable to load history right now.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_humidity', 'Humidity') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_last_seen', 'Latest reading') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_local_time', 'Local time now') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_no_data', 'No data recorded in this period.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_pressure', 'Air pressure') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_temperature', 'Air temperature') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_temperature_f', 'Air temperature (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_title', 'Safecast realtime sensor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_transport', 'Movement') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_transport_air', 'Aircraft') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_transport_bike', 'Bike or scooter') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_transport_car', 'Car or van') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_transport_unknown', 'Not detected') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'live_marker_transport_walk', 'On foot') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'locate_button_tooltip', 'Center the map on my location') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'location_error', 'An error occurred while obtaining location.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'location_permission_denied', 'Location access denied.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'location_timeout', 'Location request timed out.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'location_unavailable', 'Location unavailable.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'no_results_found', 'No results found') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'processing_complete', 'Processing complete!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'processing_on_server', 'Processing markers on server...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'qr_button_tooltip', 'QR code for the link to this map area.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'radiation_dose', 'Dose Rate') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'search_error', 'Search error. Please try again.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'search_input_tooltip', 'Search for a place by typing the first few characters of its name. A list of suggestions will appear — it may take a few seconds.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'select_files', 'Please select at least one file') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'short_link_tooltip', 'Click to copy a short share link') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'sources_full', 'We thank all who share measurements.

Anonymous uploads trace quiet paths on the map.
<a href="https://safecast.org" target="_blank">Safecast</a> fosters a global archive of readings.
<a href="https://atomfast.net" target="_blank">Atomfast</a> keeps the Atomcloud alight.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> gathers insights from Radiacode.

Each contribution widens the common picture; you are warmly invited to add your own.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'sources_title', 'Data Sources') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'speed', 'Speed') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'speed_filter_tooltip_accuracy', 'Slower measurements stay closest to the ground, so pedestrian data is the most precise.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'speed_filter_tooltip_car', 'Car: 7–200 km/h for drives and mobile surveys.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'speed_filter_tooltip_intro', 'Choose which measurements appear by travel speed.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'speed_filter_tooltip_live', 'Safecast heart: realtime data streaming from safecast.org. Toggle to show or hide live readings.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'speed_filter_tooltip_ped', 'Pedestrian: below 7 km/h for walking or stationary readings.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'speed_filter_tooltip_plane', 'Plane: 200+ km/h for airborne surveys.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'speed_filter_tooltip_title', 'Speed filters') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'theme_toggle_tooltip', 'Switch between light and dark map themes.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'title', 'Safecast New Map — Radiological Map of Planet Earth') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'track_id', 'Track ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'upload_button', 'Upload [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'upload_button_tooltip', 'Add your measurement track to the map. Supported formats: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log, .n42, .spe. You can upload multiple files, and after uploading you''ll be taken to the track page.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'upload_error', 'Error') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'waiting_for_server', 'Waiting for server...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('en', 'your_location', 'Your location') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_archive_desc', 'Descarga un paquete tgz con todos los archivos .json publicados cuando el archivo JSON está activado.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_archive_link', 'Descargar archivo semanal') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_archive_note', 'Si el archivo está desactivado, el servidor responde con HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_archive_title', 'Paquete de archivo semanal') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_latest_desc', 'Devuelve los puntos más recientes cerca de la latitud, longitud y radio solicitados en metros.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_latest_link', 'Recientes cerca de Tokio') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_latest_note', 'Ajusta lat, lon y radius_m para centrarte en tu zona.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_latest_title', 'Mediciones recientes cercanas') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_root_desc', 'Enumera metadatos, estadísticas del conjunto de datos y enlaces a todos los demás endpoints.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_root_link', 'Abrir /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_root_note', 'Empieza aquí para descubrir colecciones y el estado del servidor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_root_title', 'Índice legible por máquinas') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_track_index_desc', 'Resuelve la ruta número N y devuelve el mismo JSON que /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_track_index_link', 'Obtener índice 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_track_index_note', 'Sustituye el número para traer otra entrada. Las respuestas se transmiten como JSON delimitado por líneas.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_track_index_title', 'Descargar ruta por índice') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_tracks_desc', 'Transmite un catálogo de rutas publicadas con nombres, conteos y enlaces de descarga.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_tracks_link', 'Listar primeras rutas') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_tracks_note', 'Usa los parámetros limit y offset para paginar listas largas.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_example_tracks_title', 'Resumen de todas las rutas') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_examples_heading', 'Endpoints destacados') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_examples_note', 'Todos los enlaces responden en JSON. Cuando los endpoints transmiten datos, el navegador puede descargar archivos grandes.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_intro', 'La API JSON refleja los datos que ves en el mapa. Cada ejemplo se abre en una pestaña nueva para que revises la respuesta sin procesar.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_more_docs', '¿Necesitas más detalles?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_more_docs_link_label', 'Abrir guía detallada') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'api_title', 'Inicio rápido de la API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'attr_legal', 'Aviso legal') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'attr_license', 'Licencia') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'attr_sources', 'Fuentes de datos') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'back_to_all_tracks', 'Volver al mapa combinado de trayectorias.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'count_rate', 'Tasa de conteo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'date_slider_tooltip', 'Filtra las mediciones por fecha o rango de años.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'date_time', 'Fecha y hora') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'description', 'Mapa de radiación de Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'download_track_cim', 'Descargar recorrido (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'duration_days', '[[count]] días') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'duration_hours', '[[count]] h') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'duration_minutes', '[[count]] min') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'duration_months', '[[count]] meses') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'duration_weeks', '[[count]] semanas') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'error_during_upload', '¡Error durante la subida!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'error_processing_files', '¡Error al procesar los archivos!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'files_uploaded', 'Archivos subidos') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'geolocation_not_supported', 'Tu navegador no admite geolocalización.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'github_link_tooltip_desc', 'Proyecto de código abierto mantenido por la comunidad.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'github_link_tooltip_title', 'Safecast New Map en GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'github_link_tooltip_version', 'Versión actual: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'home_explore_global', 'Saltar y explorar el mapa mundial') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'home_or', 'o') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'home_search_placeholder', 'Introduzca una ciudad, región o país...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'home_subtitle', 'Introduzca una ubicación para comenzar.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'home_title', 'Mapa radiológico Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'home_use_location', 'Usar mi ubicación') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'legal_contact', 'Si quieres dejarnos tus comentarios, escríbenos a:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'legal_full', 'Preámbulo. Estamos creando un mapa abierto en el que personas de todo el planeta comparten lecturas de dosímetros para el bien común: ciencia, medio ambiente, educación y seguridad. Al publicar tus datos ayudas a mucha gente. Te pedimos tratar este trabajo compartido con cariño.

1) Responsabilidad. Lo que envías sigue siendo tuyo: tú respondes por la precisión y el contenido. Los datos se publican y se usan “tal cual”. El servicio no certifica su corrección, integridad ni idoneidad para ningún fin y no se hace responsable de los resultados de su uso.

2) Apertura y licencia. Al compartir mediciones, fechas, ubicaciones aproximadas, modelos de dispositivo u otros detalles, aceptas que se harán públicos y podrán emplearse libremente bajo la licencia de datos CC0 1.0 (Dominio público). El código se distribuye bajo licencia MIT. Conservas la autoría; no ofrecemos contraprestación económica; la redistribución posterior por terceros queda fuera de nuestro alcance.

3) “Tal cual” y sin verificación. El contenido se publica sin revisión previa. No podemos garantizar la calibración de los instrumentos ni la ausencia de errores. La información se comparte con fines de investigación y no constituye una recomendación profesional.

4) Privacidad y moderación. Para mantener la seguridad y la confianza, el mapa puede generalizar horas y coordenadas, y suprimir o anonimizar metadatos técnicos. También podemos ocultar o eliminar material que, a nuestro criterio razonable, sea spam, falsificado, ilegal o perturbador. Cuidamos las mediciones honestas y procuramos preservarlas.

5) Cookies. El sitio utiliza únicamente una cookie técnica de sesión de corta duración; desaparece al terminar tu visita. No guardamos más rastros.

Amigas y amigos, este mapa es fruto de un esfuerzo colectivo y de corazones abiertos. Tómenlo como un boceto del terreno, no como un plano exacto. Si nuestro trabajo te inspira, súmate: juntos podremos mejorarlo.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'legal_title', 'Aviso legal') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'legend_attention', 'Atención') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'legend_button_tooltip', 'Abrir la leyenda de niveles de radiación.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'legend_danger', 'Peligro') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'legend_full_es', 'Esta escala indica cuán seguro es un lugar para vivir, beber y comer.
Recuerda: las mediciones pueden estar incompletas; úsala solo como guía.

Verde (0–11 µR/h)
Fondo natural.
• Agua de pozo segura.
• Se puede cultivar sin pruebas.

Amarillo (11–30 µR/h)
Fondo elevado.
• Revisa agua y suelo.
• Analiza cualquier alimento antes de comer.

Rojo (30–100 µR/h)
Contaminación grave.
• No bebas el agua.
• Cultivar o comer de aquí es arriesgado; requiere pruebas de laboratorio.

Negro (>100 µR/h)
Zona crítica.
• Agua y comida inutilizables.
• Sólo estancias breves con protección.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'legend_safe', 'Seguro') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'legend_title', 'Leyenda') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'license_full', 'Este proyecto se desarrolla bajo la <a href="/LICENSE" target="_blank">MIT License</a>. El texto completo está en la raíz del repositorio y en el sitio de MIT. Puedes estudiar, compartir y modificar el código, siempre que estas libertades acompañen tu trabajo. Los datos de investigación se publican con la licencia Creative Commons 1.0 para que las mediciones sigan en el dominio público.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'license_title', 'Licencia') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_chart_all', 'Todas las lecturas') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_chart_averaged', 'Promedio en [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_chart_close', 'Cerrar') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_chart_day', 'Últimas 24 horas') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_chart_month', 'Últimos 30 días') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_chart_link', 'Abrir gráficos de radiación') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_country', 'País') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_desc_generic', 'Este sensor Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_desc_location', 'informa desde [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_desc_location_unknown', 'informa desde una zona desconocida') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_desc_named', 'Dispositivo [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_desc_no_tube', 'vigila los niveles de radiación.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_desc_radiation_sensor', 'Este es un sensor de radiación.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_desc_transport_air', 'en vuelo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_desc_transport_bike', 'en bicicleta') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_desc_transport_car', 'en coche') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_desc_transport_unknown', 'estacionario') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_desc_transport_walk', 'a pie') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_desc_tube', 'con detector [[tube]].') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_detector', 'Detector') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_device_id', 'ID del dispositivo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_device_label', 'Dispositivo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_extra_intro', 'Entorno') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_history_error', 'No se pudo cargar el historial en este momento.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_humidity', 'Humedad') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_last_seen', 'Última lectura') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_local_time', 'Hora local actual') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_no_data', 'No hay datos registrados en este período.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_pressure', 'Presión atmosférica') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_temperature', 'Temperatura del aire') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_temperature_f', 'Temperatura del aire (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_title', 'Sensor Safecast en tiempo real') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_transport', 'Movimiento') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_transport_air', 'Aeronave') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_transport_bike', 'Bicicleta o patinete') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_transport_car', 'Coche o furgoneta') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_transport_unknown', 'No detectado') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'live_marker_transport_walk', 'A pie') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'locate_button_tooltip', 'Centrar el mapa en mi ubicación') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'location_error', 'Ocurrió un error al obtener la ubicación.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'location_permission_denied', 'Acceso a la ubicación denegado.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'location_timeout', 'La solicitud de ubicación agotó el tiempo.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'location_unavailable', 'Ubicación no disponible.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'no_results_found', 'No se encontraron resultados') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'processing_complete', '¡Procesamiento completo!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'processing_on_server', 'Procesando en el servidor...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'qr_button_tooltip', 'Código QR del enlace a esta zona del mapa.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'radiation_dose', 'Tasa de dosis') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'search_error', 'Error de búsqueda. Inténtelo de nuevo.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'search_input_tooltip', 'Busque un lugar escribiendo las primeras letras de su nombre. Aparecerá una lista de sugerencias.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'select_files', 'Por favor, seleccione al menos un archivo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'short_link_tooltip', 'Haz clic para copiar un enlace corto para compartir') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'sources_full', 'Agradecemos a todas las personas que comparten mediciones.

Los envíos anónimos dibujan sendas silenciosas en el mapa.
<a href="https://safecast.org" target="_blank">Safecast</a> mantiene un archivo global de lecturas.
<a href="https://atomfast.net" target="_blank">Atomfast</a> mantiene vivo Atomcloud.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> reúne el conocimiento de Radiacode.

Cada aporte amplía la imagen común; te invitamos con cariño a sumar el tuyo.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'sources_title', 'Fuentes de datos') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'speed', 'Velocidad') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'speed_filter_tooltip_accuracy', 'Las mediciones más lentas se mantienen más cerca del suelo, por eso los datos de peatón son los más precisos.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'speed_filter_tooltip_car', 'Auto: 7–200 km/h para recorridos y mediciones móviles.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'speed_filter_tooltip_intro', 'Elige qué mediciones se muestran según la velocidad de desplazamiento.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'speed_filter_tooltip_live', 'Corazón de Safecast: datos en tiempo real de safecast.org. Usa la casilla para mostrar u ocultar las lecturas en vivo.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'speed_filter_tooltip_ped', 'Peatón: menos de 7 km/h para lecturas caminando o en reposo.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'speed_filter_tooltip_plane', 'Avión: 200+ km/h para mediciones aéreas.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'speed_filter_tooltip_title', 'Filtros de velocidad') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'theme_toggle_tooltip', 'Alterna entre el tema claro y oscuro del mapa.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'title', 'Mapa de Isótopos de Safecast — Mapa radiológico del planeta Tierra') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'track_id', 'ID de trayectoria') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'upload_button', 'Subir [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'upload_button_tooltip', 'Agregue su recorrido de mediciones al mapa. Formatos compatibles: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Puede subir varios archivos y, después de la carga, se abrirá la página del recorrido.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'upload_error', 'Error') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'waiting_for_server', 'Esperando al servidor...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('es', 'your_location', 'Tu ubicación') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_archive_desc', 'وقتی آرشیو JSON فعال باشد یک بسته tgz شامل تمام فایل‌های .json منتشرشده را دانلود می‌کند.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_archive_link', 'دانلود آرشیو هفتگی') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_archive_note', 'اگر آرشیو غیرفعال باشد، سرور خطای HTTP 404 Not Found برمی‌گرداند.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_archive_title', 'بسته آرشیو هفتگی') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_latest_desc', 'جدیدترین نقاط نزدیک به عرض، طول و شعاع درخواستی (به متر) را برمی‌گرداند.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_latest_link', 'جدیدترین‌ها نزدیک توکیو') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_latest_note', 'مقادیر lat، lon و radius_m را برای ناحیه خود تنظیم کنید.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_latest_title', 'آخرین اندازه‌گیری‌های نزدیک') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_root_desc', 'فهرست فراداده، آمار مجموعه‌داده و پیوند به سایر نقاط پایانی را نمایش می‌دهد.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_root_link', 'باز کردن /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_root_note', 'از اینجا برای کشف مجموعه‌ها و وضعیت سرور شروع کنید.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_root_title', 'نمایه قابل خواندن توسط ماشین') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_track_index_desc', 'مسیر شماره N را واکشی کرده و همان JSON /api/track/{id}.json را باز می‌گرداند.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_track_index_link', 'دریافت شماره 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_track_index_note', 'برای دریافت ورودی دیگر عدد را تغییر دهید. پاسخ‌ها به صورت JSON خط‌به‌خط ارسال می‌شوند.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_track_index_title', 'دانلود مسیر بر اساس شماره') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_tracks_desc', 'فهرستی از مسیرهای منتشرشده همراه با نام‌ها، تعداد و پیوندهای دانلود را به صورت جریان ارسال می‌کند.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_tracks_link', 'نمایش مسیرهای اول') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_tracks_note', 'برای صفحه‌بندی فهرست‌های طولانی از limit و offset استفاده کنید.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_example_tracks_title', 'خلاصه تمام مسیرها') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_examples_heading', 'نقاط پایانی منتخب') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_examples_note', 'همه پیوندها با JSON پاسخ می‌دهند. هنگام پخش داده ممکن است مرورگر فایل‌های بزرگی را دانلود کند.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_intro', 'رابط JSON همان داده‌های نقشه را برمی‌گرداند. هر نمونه در زبانه‌ای جدید باز می‌شود تا بتوانید پاسخ خام را بررسی کنید.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_more_docs', 'به راهنمای جامع‌تری نیاز دارید؟') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_more_docs_link_label', 'باز کردن راهنمای مفصل') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'api_title', 'شروع سریع API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'attr_api', 'رابط API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'attr_legal', 'اطلاعات حقوقی') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'attr_license', 'مجوز') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'attr_sources', 'منابع داده') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'back_to_all_tracks', 'بازگشت به نقشهٔ ترکیبی مسیرها.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'count_rate', 'نرخ شمارش') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'date_slider_tooltip', 'اندازه‌گیری‌ها را بر اساس تاریخ یا بازهٔ سال فیلتر کنید.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'date_time', 'تاریخ و زمان') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'description', 'نقشهٔ تابشِ چیچا') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'download_track_cim', 'دانلود مسیر (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'duration_days', '[[count]] روز') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'duration_hours', '[[count]] ساعت') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'duration_minutes', '[[count]] دقیقه') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'duration_months', '[[count]] ماه') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'duration_weeks', '[[count]] هفته') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'error_during_upload', 'خطا هنگام بارگذاری!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'error_processing_files', 'خطا در پردازش فایل‌ها!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'files_uploaded', 'فایل‌ها بارگذاری شدند') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'geolocation_not_supported', 'مرورگر شما از موقعیت‌یابی جغرافیایی پشتیبانی نمی‌کند.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'github_link_tooltip_desc', 'پروژهٔ متن‌باز با پشتیبانی جامعه.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'github_link_tooltip_title', 'نقشهٔ ایزوتوپ چیچا در GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'github_link_tooltip_version', 'نسخهٔ کنونی: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'home_explore_global', 'رد شدن و کاوش نقشه جهانی') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'home_or', 'یا') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'home_search_placeholder', 'شهر، منطقه یا کشور را وارد کنید...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'home_subtitle', 'برای شروع مکانی وارد کنید.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'home_title', 'نقشه رادیولوژیکی Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'home_use_location', 'استفاده از موقعیت من') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'legal_contact', 'برای ارسال بازخورد به این نشانی ایمیل بزنید:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'legal_full', 'پیشگفتار. ما در حال ساختن نقشه‌ای باز هستیم که مردم سراسر جهان در آن قرائت‌های دوزیمتر خود را برای خیر همگانی — برای علم، محیط زیست، آموزش و ایمنی — به اشتراک می‌گذارند. با انتشار داده‌های خود به بسیاری کمک می‌کنید. لطفاً با مهربانی از این تلاش جمعی مراقبت کنید.

۱) مسئولیت. دقت و محتوای اطلاعاتی که ارسال می‌کنید بر عهده خود شماست. داده‌ها «همان‌گونه که هستند» منتشر و استفاده می‌شوند. سرویس درستی، کامل بودن یا مناسب بودن آن‌ها را برای هیچ منظوری تأیید یا تضمین نمی‌کند و مسئول پیامدهای احتمالی استفاده از آن‌ها نیست.

۲) گشودگی و مجوز. با به اشتراک گذاشتن اندازه‌گیری‌ها، تاریخ‌ها، موقعیت‌های تقریبی، مدل دستگاه یا سایر جزئیات، می‌پذیرید که این اطلاعات در دسترس همه قرار می‌گیرد و می‌تواند تحت مجوز CC0 1.0 (قلمرو عمومی) برای داده‌ها آزادانه استفاده شود. کد پروژه همچنان تحت مجوز MIT در دسترس است. حق مؤلف برای شما محفوظ است؛ پرداختی در کار نیست؛ بازنشر بعدی توسط اشخاص ثالث خارج از کنترل ماست.

۳) «همان‌گونه که هست» و بدون راستی‌آزمایی. محتوا بدون بررسی قبلی منتشر می‌شود. ما نمی‌توانیم تضمین کنیم که دستگاه‌ها دقیقاً کالیبره شده‌اند یا عاری از خطا هستند. اطلاعات در اختیار پژوهش قرار می‌گیرد و به‌منزله توصیه حرفه‌ای نیست.

۴) حریم خصوصی و تعدیل محتوا. برای پاسداشت امنیت و اعتماد، ممکن است زمان‌ها و مختصات را کلی‌سازی کنیم و فراداده‌های فنی را حذف یا ناشناس سازیم. این حق را داریم که موادی را که به تشخیص منطقی ما اسپم، جعلی، خلاف قانون یا مختل‌کننده سرویس است پنهان یا حذف کنیم. با اندازه‌گیری‌های حسن‌نیت با دقت برخورد می‌کنیم و در حفظ آن‌ها می‌کوشیم.

۵) کوکی‌ها. این وب‌سایت تنها از یک کوکی فنی کوتاه‌مدت برای نشست استفاده می‌کند؛ با پایان بازدید شما از بین می‌رود. رد دیگری نگه نمی‌داریم.

دوستان عزیز، این نقشه حاصل تلاش مشترک و دل‌های گشوده است. آن را طرحی از چشم‌انداز بدانید، نه نقشه‌ای کاملاً دقیق. اگر این کار به دل شما می‌نشیند، کنار ما باشید — با هم می‌توانیم آن را بهتر کنیم.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'legal_title', 'اطلاعات حقوقی') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'legend_attention', 'توجه') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'legend_button_tooltip', 'راهنمای سطوح پرتو را باز کنید.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'legend_danger', 'خطر') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'legend_full_fa', 'این مقیاس نشان می‌دهد یک محل برای زندگی، آب و غذا چقدر امن است.
به یاد داشته باشید: اندازه‌گیری‌ها ممکن است ناقص باشد؛ فقط به عنوان راهنما استفاده کنید.

سبز (0–11 µR/h)
زمینه نزدیک به طبیعی.
• آب چاه معمولاً امن است.
• می‌توان بدون آزمایش گیاه کاشت.

زرد (11–30 µR/h)
زمینه بالا رفته.
• آب و خاک را بررسی کنید.
• هر غذایی را پیش از خوردن آزمایش کنید.

قرمز (30–100 µR/h)
آلودگی شدید.
• آب را ننوشید.
• کاشت یا خوردن اینجا خطرناک است؛ آزمایشگاه لازم است.

سیاه (>100 µR/h)
منطقه بحرانی.
• آب و غذا قابل استفاده نیست.
• فقط توقف کوتاه با حفاظت.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'legend_safe', 'ایمن') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'legend_title', 'راهنما') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'license_full', 'این پروژه تحت مجوز <a href="/LICENSE" target="_blank">MIT License</a> رشد می‌کند. متن کامل در ریشه مخزن و در وب‌سایت MIT در دسترس است. می‌توانید کد را مطالعه، به اشتراک بگذارید و تغییر دهید؛ به شرطی که این آزادی‌ها همراه کار شما بماند. داده‌های پژوهشی با مجوز Creative Commons 1.0 منتشر می‌شوند تا اندازه‌گیری‌ها در حوزه عمومی باقی بمانند.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'license_title', 'مجوز') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_chart_all', 'همه خوانش‌ها') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_chart_averaged', 'میانگین در [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_chart_close', 'بستن') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_chart_day', '۲۴ ساعت گذشته') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_chart_month', '۳۰ روز گذشته') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_chart_link', 'باز کردن نمودارهای تابش') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_country', 'کشور') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_desc_generic', 'این حسگر Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_desc_location', 'از [[place]] گزارش می‌دهد') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_desc_location_unknown', 'از منطقه‌ای ناشناخته گزارش می‌دهد') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_desc_named', 'دستگاه [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_desc_no_tube', 'سطح تابش را زیر نظر دارد.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_desc_radiation_sensor', 'این یک حسگر تابش است.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_desc_transport_air', 'در حال پرواز') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_desc_transport_bike', 'با دوچرخه') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_desc_transport_car', 'با خودرو') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_desc_transport_unknown', 'ثابت') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_desc_transport_walk', 'پیاده') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_desc_tube', 'با آشکارساز [[tube]].') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_detector', 'آشکارساز') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_device_id', 'شناسه دستگاه') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_device_label', 'دستگاه') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_extra_intro', 'محیط') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_history_error', 'بارگذاری تاریخچه در حال حاضر امکان‌پذیر نیست.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_humidity', 'رطوبت') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_last_seen', 'آخرین خوانش') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_local_time', 'ساعت محلی') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_no_data', 'داده‌ای در این دوره ثبت نشده است.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_pressure', 'فشار هوا') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_temperature', 'دمای هوا') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_temperature_f', 'دمای هوا (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_title', 'حسگر بلادرنگ Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_transport', 'حرکت') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_transport_air', 'هواپیما') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_transport_bike', 'دوچرخه یا اسکوتر') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_transport_car', 'خودرو یا ون') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_transport_unknown', 'تشخیص داده نشد') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'live_marker_transport_walk', 'پیاده') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'locate_button_tooltip', 'نقشه را روی موقعیت من متمرکز کن') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'location_error', 'هنگام دریافت موقعیت خطایی رخ داد.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'location_permission_denied', 'دسترسی به موقعیت رد شد.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'location_timeout', 'درخواست موقعیت زمان‌بر شد.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'location_unavailable', 'موقعیت در دسترس نیست.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'no_results_found', 'نتیجه‌ای یافت نشد') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'processing_complete', 'پردازش انجام شد!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'processing_on_server', 'در حال پردازش روی سرور...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'qr_button_tooltip', 'کد QRِ لینک به این بخش از نقشه.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'radiation_dose', 'نرخ دوز') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'search_error', 'خطا در جستجو. لطفاً دوباره تلاش کنید.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'search_input_tooltip', 'با تایپ حروف اول نام مکان جستجو کنید. فهرستی از پیشنهادها نمایش داده می‌شود.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'select_files', 'لطفاً دست‌کم یک فایل انتخاب کنید') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'short_link_tooltip', 'برای کپی کردن لینک کوتاه اشتراک کلیک کنید') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'sources_full', 'از همه کسانی که اندازه‌گیری‌ها را به اشتراک می‌گذارند سپاسگزاریم.

بارگذاری‌های ناشناس ردپاهای آرامی روی نقشه ترسیم می‌کنند.
<a href="https://safecast.org" target="_blank">Safecast</a> آرشیوی جهانی از قرائت‌ها را نگه می‌دارد.
<a href="https://atomfast.net" target="_blank">Atomfast</a> ابر Atomcloud را روشن نگه می‌دارد.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> بینش‌های Radiacode را گرد می‌آورد.

هر مشارکت تصویر مشترک را گسترده‌تر می‌کند؛ با آغوش باز از شما دعوت می‌کنیم سهم خود را بیفزایید.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'sources_title', 'منابع داده') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'speed', 'سرعت') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'speed_filter_tooltip_accuracy', 'اندازه‌گیری‌های آهسته‌تر نزدیک‌تر به سطح زمین می‌مانند، بنابراین داده‌های پیاده دقیق‌ترین هستند.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'speed_filter_tooltip_car', 'خودرو: ۷ تا ۲۰۰ کیلومتر بر ساعت برای رانندگی و پایش سیار.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'speed_filter_tooltip_intro', 'انتخاب کنید کدام اندازه‌گیری‌ها بر اساس سرعت حرکت نمایش داده شوند.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'speed_filter_tooltip_live', 'قلب Safecast: داده‌های لحظه‌ای از safecast.org. با فعال‌سازی این گزینه می‌توانید قرائت‌های زنده را نشان دهید یا پنهان کنید.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'speed_filter_tooltip_ped', 'پیاده: کمتر از ۷ کیلومتر بر ساعت برای اندازه‌گیری در حالت پیاده یا ثابت.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'speed_filter_tooltip_plane', 'هواپیما: +۲۰۰ کیلومتر بر ساعت برای پایش هوایی.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'speed_filter_tooltip_title', 'فیلترهای سرعت') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'theme_toggle_tooltip', 'بین تم روشن و تیرهٔ نقشه جابه‌جا شوید.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'title', 'نقشه ایزوتوپ‌های چیچا — نقشهٔ رادیولوژیِ سیارهٔ زمین') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'track_id', 'شناسهٔ مسیر') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'upload_button', 'بارگذاری [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'upload_button_tooltip', 'مسیر اندازه‌گیری خود را به نقشه اضافه کنید. قالب‌های پشتیبانی‌شده: ‎.kml, .kmz, .gpx, .csv, .rctrk, .json, .log‎. می‌توانید چندین فایل بارگذاری کنید و پس از بارگذاری، صفحهٔ مسیر باز می‌شود.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'upload_error', 'خطا') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'waiting_for_server', 'در انتظار سرور...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fa', 'your_location', 'موقعیت شما') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_archive_desc', 'Lataa tgz-paketin, jossa on kaikki julkaistut .json-tiedostot, kun JSON-arkisto on käytössä.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_archive_link', 'Lataa viikkoarkisto') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_archive_note', 'Jos arkisto on pois päältä, palvelin palauttaa HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_archive_title', 'Viikoittainen arkivipaketti') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_latest_desc', 'Palauttaa tuoreimmat pisteet lähellä annettua leveyttä, pituutta ja hakusädettä metreissä.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_latest_link', 'Viimeisimmät Tokion lähellä') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_latest_note', 'Säädä lat-, lon- ja radius_m-arvoja omaa aluettasi varten.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_latest_title', 'Viimeisimmät mittaukset lähistöllä') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_root_desc', 'Listaa metatiedot, aineiston tilastot ja linkit jokaiseen muuhun päätepisteeseen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_root_link', 'Avaa /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_root_note', 'Aloita täältä kokoelmien ja palvelimen tilan tutkiminen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_root_title', 'Koneluettava indeksi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_track_index_desc', 'Hakee N:nnen reitin ja palauttaa saman JSONin kuin /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_track_index_link', 'Nouda indeksi 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_track_index_note', 'Vaihda numero hakeaksesi toisen kohteen. Vastaukset striimataan rivitettynä JSONina.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_track_index_title', 'Lataa reitti indeksin perusteella') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_tracks_desc', 'Striimaa julkaistujen reittien luettelon nimiineen, lukumäärineen ja latauslinkkeineen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_tracks_link', 'Luettele ensimmäiset reitit') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_tracks_note', 'Käytä limit- ja offset-parametreja pitkien listojen sivutukseen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_example_tracks_title', 'Kaikkien reittien yhteenvedot') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_examples_heading', 'Nostetut päätepisteet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_examples_note', 'Kaikki linkit palauttavat JSONia. Kun päätepisteet striimaavat dataa, selain voi ladata suuria tiedostoja.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_intro', 'JSON-rajapinta heijastaa kartalla näkyvät tiedot. Jokainen esimerkki avautuu uuteen välilehteen, jotta voit tutkia raakaa vastausta.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_more_docs', 'Tarvitsetko kattavamman ohjeen?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_more_docs_link_label', 'Avaa yksityiskohtainen opas') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'api_title', 'API-pikaopas') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'attr_legal', 'Oikeudelliset tiedot') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'attr_license', 'Lisenssi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'attr_sources', 'Tietolähteet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'back_to_all_tracks', 'Takaisin yhdistettyyn reittikarttaan.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'count_rate', 'Laskentanopeus') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'date_slider_tooltip', 'Suodata mittaukset päivämäärän tai vuosivälin mukaan.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'date_time', 'Päivämäärä ja aika') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'description', 'Safecastn säteilykartta') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'download_track_cim', 'Lataa reitti (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'duration_days', '[[count]] päivää') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'duration_hours', '[[count]] t') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'duration_minutes', '[[count]] min') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'duration_months', '[[count]] kuukautta') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'duration_weeks', '[[count]] viikkoa') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'error_during_upload', 'Virhe lähetyksen aikana!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'error_processing_files', 'Virhe tiedostoja käsiteltäessä!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'files_uploaded', 'Tiedostot lähetetty') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'geolocation_not_supported', 'Selaimesi ei tue paikkatietoja.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'github_link_tooltip_desc', 'Avoimen lähdekoodin projekti yhteisön ylläpitämänä.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'github_link_tooltip_title', 'Safecast New Map GitHubissa') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'github_link_tooltip_version', 'Nykyinen versio: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'home_explore_global', 'Ohita ja tutustu maailmankarttaan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'home_or', 'tai') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'home_search_placeholder', 'Syötä kaupunki, alue tai maa...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'home_subtitle', 'Syötä sijainti aloittaaksesi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'home_title', 'Safecast säteilykartta') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'home_use_location', 'Käytä sijaintiani') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'legal_contact', 'Jos haluat antaa palautetta, kirjoita osoitteeseen:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'legal_full', 'Esipuhe. Rakennamme avointa karttaa, jossa ihmiset kaikkialta maailmasta jakavat dosimetrihavaintoja yhteiseksi hyväksi – tieteen, ympäristön, koulutuksen ja turvallisuuden tueksi. Kun julkaiset tietosi, autat monia. Pyydämme, että kohtelet tätä yhteistä työtä lämpimästi ja kunnioittaen.

1) Vastuu. Olet itse vastuussa lähettämiesi tietojen paikkansapitävyydestä ja sisällöstä. Tiedot julkaistaan ja niitä käytetään “sellaisinaan”. Palvelu ei vahvista eikä takaa niiden oikeellisuutta, täydellisyyttä tai soveltuvuutta mihinkään tarkoitukseen eikä vastaa niiden käytön seurauksista.

2) Avoimuus ja lisenssi. Kun jaat mittauksia, päivämääriä, likimääräisiä sijainteja, laitemalleja tai muita tietoja, ymmärrät, että ne tulevat kaikkien saataville ja niitä voidaan käyttää vapaasti CC0 1.0 (Public Domain) -lisenssin mukaisesti. Ohjelmakoodi on jatkossakin MIT-lisenssin alaisena. Tekijyys säilyy sinulla; korvausta ei makseta; kolmansien osapuolten jatkolevitys ei ole hallinnassamme.

3) “Sellaisinaan” ilman ennakkotarkistusta. Julkaisut ilmestyvät ilman ennakkovalvontaa. Emme voi taata laitteiden kalibrointia tai virheettömyyttä. Tiedot jaetaan tutkimuskäyttöön, eivätkä ne ole ammatillisia suosituksia.

4) Yksityisyys ja moderointi. Turvallisuuden ja luottamuksen vuoksi voimme yleistää aikaleimoja ja koordinaatteja sekä poistaa tai anonymisoida teknisiä metatietoja. Voimme piilottaa tai poistaa aineistoa, jonka katsomme kohtuudella olevan roskaa, väärennöksiä, lainvastaista tai palvelua häiritsevää. Aidosti kerättyjä mittauksia käsittelemme huolella ja pyrimme säilyttämään ne.

5) Evästeet. Sivusto käyttää vain lyhytikäistä teknistä istuntoevästettä; se poistuu käyntisi päätyttyä. Emme tallenna muita jälkiä.

Ystävät, tämä kartta on yhteisen ponnistuksen ja avointen sydänten tulos. Kohdelkaa sitä maaston luonnoksena, ei millintarkkana piirustuksena. Jos työmme puhuttelee sinua, liity joukkoon – yhdessä teemme siitä entistä paremman.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'legal_title', 'Lakiasiat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'legend_attention', 'Huomio') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'legend_button_tooltip', 'Avaa säteilytasojen selite.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'legend_danger', 'Vaara') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'legend_full_fi', 'Tämä asteikko kertoo, kuinka turvallinen paikka on elämälle, vedelle ja ruoalle.
Muista: mittaukset voivat olla puutteellisia; käytä vain ohjeena.

Vihreä (0–11 µR/h)
Lähes luonnollinen tausta.
• Kaivovesi yleensä turvallista.
• Voit kasvattaa kasveja ilman testejä.

Keltainen (11–30 µR/h)
Kohonnut tausta.
• Tarkista vesi ja maa.
• Tutki ruoka ennen syömistä.

Punainen (30–100 µR/h)
Vakava saastuminen.
• Älä juo vettä.
• Täällä viljely tai syöminen on riskialtista; laboratoriotestit pakollisia.

Musta (>100 µR/h)
Kriittinen alue.
• Vesi ja ruoka käyttökelvottomia.
• Pitkä oleskelu mahdoton; vain lyhyet käynnit suojassa.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'legend_safe', 'Turvallinen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'legend_title', 'Selite') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'license_full', 'Tämä projekti elää <a href="/LICENSE" target="_blank">MIT License</a> -lisenssin alla. Koko teksti löytyy projektin juuresta ja MIT:n sivuilta. Voit tutkia, jakaa ja muokata koodia, kunhan viet nämä vapaudet mukanasi omiin töihisi. Tutkimusaineistot julkaistaan Creative Commons 1.0 -lisenssillä, jotta mittaukset pysyvät julkisessa omistuksessa.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'license_title', 'Lisenssi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_chart_all', 'Kaikki lukemat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_chart_averaged', 'Keskiarvo [[window]] ajalta') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_chart_close', 'Sulje') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_chart_day', 'Viimeiset 24 tuntia') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_chart_month', 'Viimeiset 30 päivää') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_chart_link', 'Avaa säteilykaaviot') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_country', 'Maa') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_desc_generic', 'Tämä Safecast-anturi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_desc_location', 'raportoi paikasta [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_desc_location_unknown', 'raportoi tuntemattomalta alueelta') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_desc_named', 'Laite [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_desc_no_tube', 'valvoo säteilytasoja.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_desc_radiation_sensor', 'Tämä on säteilyanturi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_desc_transport_air', 'lennon aikana') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_desc_transport_bike', 'pyörällä') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_desc_transport_car', 'autolla') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_desc_transport_unknown', 'paikallaan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_desc_transport_walk', 'jalan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_desc_tube', '[[tube]]-ilmaisimella.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_detector', 'Ilmaisin') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_device_id', 'Laitteen tunnus') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_device_label', 'Laite') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_extra_intro', 'Ympäristö') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_history_error', 'Historiaa ei voida ladata juuri nyt.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_humidity', 'Kosteus') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_last_seen', 'Viimeisin lukema') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_local_time', 'Paikallinen aika nyt') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_no_data', 'Ei tallennettuja tietoja tällä jaksolla.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_pressure', 'Ilmanpaine') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_temperature', 'Ilman lämpötila') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_temperature_f', 'Ilman lämpötila (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_title', 'Safecast reaaliaikainen anturi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_transport', 'Liike') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_transport_air', 'Lentokone') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_transport_bike', 'Pyörä tai potkulauta') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_transport_car', 'Auto tai pakettiauto') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_transport_unknown', 'Ei havaittu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'live_marker_transport_walk', 'Jalan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'locate_button_tooltip', 'Keskitä kartta sijaintiini') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'location_error', 'Virhe sijaintia haettaessa.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'location_permission_denied', 'Pääsy sijaintiin estetty.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'location_timeout', 'Sijaintipyynnön aikakatkaisu.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'location_unavailable', 'Sijainti ei ole käytettävissä.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'no_results_found', 'Tuloksia ei löytynyt') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'processing_complete', 'Käsittely valmis!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'processing_on_server', 'Käsitellään palvelimella...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'qr_button_tooltip', 'Tämän kartta-alueen linkin QR-koodi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'radiation_dose', 'Annosnopeus') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'search_error', 'Hakuvirhe. Yritä uudelleen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'search_input_tooltip', 'Etsi paikkaa kirjoittamalla nimen alkukirjaimet. Ehdotuslista tulee näkyviin.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'select_files', 'Valitse vähintään yksi tiedosto') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'short_link_tooltip', 'Napsauta kopioidaksesi lyhyen jakolinkin') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'sources_full', 'Kiitämme kaikkia mittauksia jakavia.

Anonyymit lataukset piirtävät hiljaisia reittejä kartalle.
<a href="https://safecast.org" target="_blank">Safecast</a> ylläpitää maailmanlaajuista mittausarkistoa.
<a href="https://atomfast.net" target="_blank">Atomfast</a> pitää Atomcloudin käynnissä.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> kokoaa Radiacode-havaintoja.

Jokainen panos laajentaa yhteistä kuvaa; olet lämpimästi tervetullut lisäämään omasi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'sources_title', 'Datalähteet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'speed', 'Nopeus') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'speed_filter_tooltip_accuracy', 'Hitaammat mittaukset pysyvät lähimpänä maanpintaa, joten jalankulkijatiedot ovat tarkimpia.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'speed_filter_tooltip_car', 'Auto: 7–200 km/h ajomittauksiin ja liikkuviin keräyksiin.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'speed_filter_tooltip_intro', 'Valitse, mitkä mittaukset näytetään liikkumisnopeuden perusteella.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'speed_filter_tooltip_live', 'Safecast-sydän: reaaliaikaiset tiedot safecast.orgista. Valitse tästä, näytetäänkö vai piilotetaanko live-mittaukset.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'speed_filter_tooltip_ped', 'Jalankulkija: alle 7 km/h kävellessä tai paikallaan mitattaessa.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'speed_filter_tooltip_plane', 'Lentokone: 200+ km/h ilmamittauksiin.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'speed_filter_tooltip_title', 'Nopeussuodattimet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'theme_toggle_tooltip', 'Vaihda kartan vaalean ja tumman teeman välillä.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'title', 'Safecastn isotooppikartta — Maapallon radiologinen kartta') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'track_id', 'Reitin tunnus') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'upload_button', 'Lähetä [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'upload_button_tooltip', 'Lisää mittausreittisi kartalle. Tuetut muodot: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Voit lähettää useita tiedostoja; lähetyksen jälkeen avautuu reittisivu.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'upload_error', 'Virhe') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'waiting_for_server', 'Odotetaan palvelinta...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fi', 'your_location', 'Sijaintisi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_archive_desc', 'Télécharge une archive tgz avec tous les fichiers .json publiés lorsque l’archive JSON est activée.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_archive_link', 'Télécharger l’archive hebdomadaire') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_archive_note', 'Si l’archive est désactivée, le serveur renvoie HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_archive_title', 'Archive hebdomadaire') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_latest_desc', 'Renvoie les points les plus récents près de la latitude, de la longitude et du rayon recherchés en mètres.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_latest_link', 'Récentes près de Tokyo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_latest_note', 'Ajustez lat, lon et radius_m pour votre zone.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_latest_title', 'Mesures récentes à proximité') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_root_desc', 'Liste les métadonnées, les statistiques du jeu de données et des liens vers chaque autre point de terminaison.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_root_link', 'Ouvrir /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_root_note', 'Commencez ici pour découvrir les collections et l’état du serveur.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_root_title', 'Index lisible par machine') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_track_index_desc', 'Résout le N-ième parcours et renvoie le même JSON que /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_track_index_link', 'Récupérer l’index 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_track_index_note', 'Remplacez le nombre pour obtenir une autre entrée. Les réponses sont diffusées en JSON ligne par ligne.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_track_index_title', 'Télécharger un parcours par index') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_tracks_desc', 'Diffuse un catalogue des parcours publiés avec noms, décomptes et liens de téléchargement.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_tracks_link', 'Lister les premiers parcours') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_tracks_note', 'Utilisez les paramètres limit et offset pour parcourir les longues listes.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_example_tracks_title', 'Tous les résumés de parcours') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_examples_heading', 'Points de terminaison mis en avant') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_examples_note', 'Tous les liens renvoient du JSON. Lorsque les points de terminaison diffusent des données, le navigateur peut télécharger de gros fichiers.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_intro', 'L’API JSON reflète les données affichées sur la carte. Chaque exemple s’ouvre dans un nouvel onglet pour inspecter la réponse brute.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_more_docs', 'Besoin d’une documentation plus complète ?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_more_docs_link_label', 'Ouvrir le guide détaillé') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'api_title', 'Démarrage rapide de l’API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'attr_legal', 'Mentions légales') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'attr_license', 'Licence') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'attr_sources', 'Sources de données') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'back_to_all_tracks', 'Retour à la carte combinée des pistes.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'count_rate', 'Taux de comptage') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'date_slider_tooltip', 'Filtrer les mesures par date ou plage d''années.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'date_time', 'Date et heure') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'description', 'Carte de radiations de Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'download_track_cim', 'Télécharger la piste (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'duration_days', '[[count]] jours') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'duration_hours', '[[count]] h') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'duration_minutes', '[[count]] min') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'duration_months', '[[count]] mois') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'duration_weeks', '[[count]] semaines') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'error_during_upload', 'Erreur lors du téléversement !') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'error_processing_files', 'Erreur lors du traitement des fichiers !') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'files_uploaded', 'Fichiers téléversés') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'geolocation_not_supported', 'Votre navigateur ne prend pas en charge la géolocalisation.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'github_link_tooltip_desc', 'Projet open source maintenu par la communauté.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'github_link_tooltip_title', 'Safecast New Map sur GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'github_link_tooltip_version', 'Version actuelle : {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'home_explore_global', 'Passer et explorer la carte mondiale') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'home_or', 'ou') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'home_search_placeholder', 'Entrez une ville, région ou pays...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'home_subtitle', 'Entrez un lieu pour commencer.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'home_title', 'Carte radiologique Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'home_use_location', 'Utiliser ma position') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'legal_contact', 'Pour toute remarque, écrivez-nous à :') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'legal_full', 'Préambule. Nous créons une carte ouverte où des personnes du monde entier partagent les mesures de leurs dosimètres pour le bien commun — science, environnement, éducation et sécurité. En publiant vos données, vous rendez service à beaucoup. Merci de prendre soin de ce projet collectif.

1) Responsabilité. Vous demeurez responsable de l’exactitude et du contenu des informations que vous transmettez. Les données sont publiées et utilisées « en l’état ». Le service n’en certifie ni l’exactitude, ni l’exhaustivité, ni l’adéquation à quelque usage que ce soit et décline toute responsabilité quant aux conséquences de leur utilisation.

2) Ouverture et licence. En partageant vos mesures, dates, localisations approximatives, modèles d’appareils ou autres détails, vous acceptez qu’ils deviennent accessibles à toutes et à tous et puissent être employés librement sous la licence de données CC0 1.0 (domaine public). Le code reste diffusé sous licence MIT. Vous conservez votre paternité ; aucune rémunération n’est prévue ; la redistribution par des tiers ne dépend pas de nous.

3) « En l’état » et sans vérification préalable. Les contributions sont publiées sans contrôle préalable. Nous ne pouvons garantir ni l’étalonnage des instruments ni l’absence d’erreurs. Les informations sont partagées à des fins de recherche et ne constituent pas un avis professionnel.

4) Confidentialité et modération. Pour préserver la sécurité et la confiance, nous pouvons généraliser horaires et coordonnées, et supprimer ou anonymiser certains métadonnées techniques. Nous nous réservons la possibilité de masquer ou d’effacer tout contenu qui, selon notre appréciation raisonnable, relèverait du spam, serait falsifié, illégal ou perturberait le service. Les mesures réalisées de bonne foi sont traitées avec soin et nous cherchons à les préserver.

5) Cookies. Le site n’utilise qu’un cookie technique de session de courte durée ; il disparaît lorsque votre visite s’achève. Aucun autre traceur n’est conservé.

Amies et amis, cette carte est le fruit d’un effort commun et de cœurs ouverts. Considérez-la comme un croquis du terrain, non comme un plan millimétré. Si notre travail vous parle, rejoignez-nous : ensemble nous la ferons progresser.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'legal_title', 'Mentions légales') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'legend_attention', 'Attention') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'legend_button_tooltip', 'Ouvrir la légende des niveaux de radiation.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'legend_danger', 'Danger') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'legend_full_fr', 'Cette échelle indique à quel point un lieu est sûr pour la vie, l’eau et la nourriture.
Souviens‑toi : les mesures peuvent être incomplètes ; n’en fais qu’un guide.

Vert (0–11 µR/h)
Fond naturel.
• L’eau de puits est généralement potable.
• On peut cultiver sans tests.

Jaune (11–30 µR/h)
Fond élevé.
• Vérifie l’eau et le sol.
• Analyse toute nourriture avant de la consommer.

Rouge (30–100 µR/h)
Contamination sérieuse.
• Ne bois pas l’eau.
• Cultiver ou manger ici est risqué ; analyses en laboratoire obligatoires.

Noir (>100 µR/h)
Zone critique.
• Eau et nourriture inutilisables.
• Séjour prolongé impossible ; visites courtes avec protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'legend_safe', 'Sûr') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'legend_title', 'Légende') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'license_full', 'Ce projet s’épanouit sous la licence <a href="/LICENSE" target="_blank">MIT License</a>. Le texte intégral se trouve à la racine du dépôt et sur le site du MIT. Vous pouvez étudier, partager et modifier le code, pourvu que ces libertés accompagnent vos travaux. Les données de recherche sont publiées sous licence Creative Commons 1.0 afin que les mesures demeurent dans le domaine public.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'license_title', 'Licence') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_chart_all', 'Toutes les lectures') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_chart_averaged', 'Moyenne sur [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_chart_close', 'Fermer') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_chart_day', 'Dernières 24 heures') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_chart_month', 'Derniers 30 jours') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_chart_link', 'Ouvrir les graphiques de radiation') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_country', 'Pays') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_desc_generic', 'Ce capteur Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_desc_location', 'rapporte depuis [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_desc_location_unknown', 'rapporte depuis une zone inconnue') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_desc_named', 'Appareil [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_desc_no_tube', 'surveille les niveaux de radiation.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_desc_radiation_sensor', 'Ceci est un capteur de radiation.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_desc_transport_air', 'en vol') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_desc_transport_bike', 'à vélo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_desc_transport_car', 'en voiture') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_desc_transport_unknown', 'stationnaire') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_desc_transport_walk', 'à pied') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_desc_tube', 'avec le détecteur [[tube]].') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_detector', 'Détecteur') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_device_id', 'ID de l''appareil') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_device_label', 'Appareil') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_extra_intro', 'Environnement') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_history_error', 'Impossible de charger l''historique pour le moment.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_humidity', 'Humidité') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_last_seen', 'Dernière lecture') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_local_time', 'Heure locale actuelle') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_no_data', 'Aucune donnée enregistrée pour cette période.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_pressure', 'Pression atmosphérique') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_temperature', 'Température de l''air') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_temperature_f', 'Température de l''air (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_title', 'Capteur Safecast en temps réel') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_transport', 'Mouvement') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_transport_air', 'Avion') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_transport_bike', 'Vélo ou trottinette') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_transport_car', 'Voiture ou camionnette') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_transport_unknown', 'Non détecté') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'live_marker_transport_walk', 'À pied') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'locate_button_tooltip', 'Centrer la carte sur ma position') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'location_error', 'Une erreur est survenue lors de l’obtention de la localisation.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'location_permission_denied', 'Accès à la localisation refusé.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'location_timeout', 'La demande de localisation a expiré.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'location_unavailable', 'Localisation indisponible.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'no_results_found', 'Aucun résultat trouvé') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'processing_complete', 'Traitement terminé !') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'processing_on_server', 'Traitement sur le serveur...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'qr_button_tooltip', 'Code QR du lien vers cette zone de la carte.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'radiation_dose', 'Débit de dose') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'search_error', 'Erreur de recherche. Veuillez réessayer.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'search_input_tooltip', 'Recherchez un lieu en tapant les premières lettres de son nom. Une liste de suggestions apparaîtra.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'select_files', 'Veuillez sélectionner au moins un fichier') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'short_link_tooltip', 'Cliquez pour copier un lien court à partager') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'sources_full', 'Nous remercions toutes celles et ceux qui partagent leurs mesures.

Les envois anonymes tracent des chemins discrets sur la carte.
<a href="https://safecast.org" target="_blank">Safecast</a> entretient une archive mondiale des relevés.
<a href="https://atomfast.net" target="_blank">Atomfast</a> maintient l’Atomcloud allumée.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> recueille les connaissances issues de Radiacode.

Chaque contribution élargit la vision commune ; nous vous invitons chaleureusement à apporter la vôtre.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'sources_title', 'Sources de données') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'speed', 'Vitesse') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'speed_filter_tooltip_accuracy', 'Les mesures plus lentes restent proches du sol ; les données à pied sont donc les plus précises.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'speed_filter_tooltip_car', 'Voiture : 7 à 200 km/h pour les trajets et relevés mobiles.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'speed_filter_tooltip_intro', 'Choisissez quelles mesures afficher selon la vitesse de déplacement.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'speed_filter_tooltip_live', 'Cœur Safecast : données en temps réel depuis safecast.org. Utilisez la case pour afficher ou masquer les relevés en direct.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'speed_filter_tooltip_ped', 'Piéton : moins de 7 km/h pour les mesures à pied ou à l''arrêt.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'speed_filter_tooltip_plane', 'Avion : 200+ km/h pour les relevés aériens.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'speed_filter_tooltip_title', 'Filtres de vitesse') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'theme_toggle_tooltip', 'Basculer entre les thèmes clair et sombre de la carte.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'title', 'Carte des isotopes de Safecast — Carte radiologique de la planète Terre') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'track_id', 'ID de piste') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'upload_button', 'Téléverser [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'upload_button_tooltip', 'Ajoutez votre parcours de mesures à la carte. Formats pris en charge : .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Vous pouvez téléverser plusieurs fichiers et, après le téléversement, la page du parcours s’ouvrira.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'upload_error', 'Erreur') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'waiting_for_server', 'En attente du serveur...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('fr', 'your_location', 'Votre position') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_archive_desc', 'מורידה חבילת tgz עם כל קבצי .json שפורסמו כאשר ארכיון ה-JSON מופעל.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_archive_link', 'הורידו ארכיון שבועי') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_archive_note', 'אם הארכיון כבוי השרת מחזיר HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_archive_title', 'חבילת ארכיון שבועית') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_latest_desc', 'מחזיר את הסמנים העדכניים ביותר ליד קו הרוחב, קו האורך והרדיוס המבוקשים במטרים.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_latest_link', 'עדכניות ליד טוקיו') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_latest_note', 'התאימו את lat, lon ו-radius_m לאזור שלכם.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_latest_title', 'מדידות עדכניות בקרבת מקום') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_root_desc', 'מציג מטא-נתונים, סטטיסטיקות של מערך הנתונים וקישורים לכל נקודת קצה נוספת.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_root_link', 'פתחו /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_root_note', 'התחילו כאן כדי לגלות אוספים ומצב שרת.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_root_title', 'אינדקס קריא למכונה') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_track_index_desc', 'מאחזר את המסלול ה-N ומחזיר את אותו JSON כמו /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_track_index_link', 'שלפו אינדקס 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_track_index_note', 'החליפו את המספר כדי להביא רשומה אחרת. התגובות מוזרמות כ-JSON לפי שורה.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_track_index_title', 'הורדת מסלול לפי אינדקס') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_tracks_desc', 'מזרים קטלוג של מסלולים שפורסמו כולל שמות, מונים וקישורי הורדה.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_tracks_link', 'הצגת המסלולים הראשונים') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_tracks_note', 'השתמשו בפרמטרים limit ו-offset כדי לדפדף ברשימות ארוכות.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_example_tracks_title', 'סיכומי כל המסלולים') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_examples_heading', 'נקודות קצה מומלצות') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_examples_note', 'כל הקישורים מחזירים JSON. כאשר נקודות הקצה מזריםות נתונים, הדפדפן עשוי להוריד קבצים גדולים.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_intro', 'ממשק ה-JSON משקף את הנתונים שמוצגים במפה. כל דוגמה נפתחת בלשונית חדשה כדי שתוכלו לבדוק את התשובה הגולמית.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_more_docs', 'צריכים תיעוד מעמיק יותר?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_more_docs_link_label', 'פתחו מדריך מפורט') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'api_title', 'פתיחה מהירה ל-API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'attr_api', 'ממשק API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'attr_legal', 'מידע משפטי') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'attr_license', 'רישיון') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'attr_sources', 'מקורות נתונים') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'back_to_all_tracks', 'חזרה למפת המסלולים המשולבת.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'count_rate', 'קצב ספירה') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'date_slider_tooltip', 'סננו מדידות לפי תאריך או טווח שנים.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'date_time', 'תאריך ושעה') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'description', 'מפת הקרינה של Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'download_track_cim', 'הורדת מסלול (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'duration_days', '[[count]] ימים') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'duration_hours', '[[count]] שע׳') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'duration_minutes', '[[count]] דק׳') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'duration_months', '[[count]] חודשים') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'duration_weeks', '[[count]] שבועות') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'error_during_upload', 'שגיאה במהלך ההעלאה!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'error_processing_files', 'שגיאה בעיבוד הקבצים!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'files_uploaded', 'קבצים הועלו') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'geolocation_not_supported', 'הדפדפן שלך אינו תומך בקביעת מיקום.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'github_link_tooltip_desc', 'פרויקט קוד פתוח שמנוהל על ידי הקהילה.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'github_link_tooltip_title', 'Safecast New Map ב-GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'github_link_tooltip_version', 'גרסה נוכחית: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'home_explore_global', 'דלג וחקור את המפה העולמית') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'home_or', 'או') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'home_search_placeholder', 'הזן עיר, אזור או מדינה...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'home_subtitle', 'הזן מיקום כדי להתחיל.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'home_title', 'מפת קרינה של Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'home_use_location', 'השתמש במיקום שלי') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'legal_contact', 'לשיתוף משוב כתבו אל:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'legal_full', 'פתיח. אנחנו בונים מפה פתוחה שבה אנשים מכל רחבי העולם משתפים את קריאות הדוזימטר לטובת הכלל — למען המדע, הסביבה, החינוך והבטיחות. כשאתם מפרסמים נתונים אתם מסייעים לרבים. אנא התייחסו אל המיזם המשותף הזה ברגישות ובכבוד.

1) אחריות. הדיוק והתוכן של המידע שאתם שולחים באחריותכם. הנתונים מתפרסמים ומשמשים “כפי שהם”. השירות אינו מאשר ואינו מבטיח את נכונותם, שלמותם או התאמתם לכל מטרה שהיא, ואינו אחראי לתוצאות אפשריות של השימוש בהם.

2) פתיחות ורישוי. כשאתם משתפים מדידות, תאריכים, מיקום משוער, דגם מכשיר או פרטים נוספים, אתם מבינים שהמידע יהיה פתוח לכול ויוכל לשמש בחופשיות תחת רישיון הנתונים CC0 1.0 (נחלת הכלל). הקוד נשאר מופץ תחת רישיון MIT. הזכות ליצירה נשמרת לכם; אין תשלום; הפצה נוספת בידי צדדים שלישיים אינה בשליטתנו.

3) “כפי שהם” וללא אימות מוקדם. הפרסומים מופיעים ללא בדיקה מוקדמת. איננו יכולים להבטיח כיול מכשירים או היעדר טעויות. המידע נמסר לצורכי מחקר ואינו מהווה ייעוץ מקצועי.

4) פרטיות וניהול. כדי לשמור על ביטחון ואמון אנו עשויים לעגל זמני מדידה וקואורדינטות, ולמחוק או לאנונם נתוני מטא טכניים. אנו רשאים להסתיר או למחוק חומרים שלפי שיקול דעתנו הסביר הם ספאם, זיוף, עבירה על החוק או פגיעה בשירות. למדידות הנעשות בתום לב אנו מתייחסים בזהירות ומשתדלים לשמור עליהן.

5) עוגיות. האתר עושה שימוש בעוגיית סשן טכנית וקצרת-מועד בלבד; היא נמחקת בסיום הביקור. איננו שומרים עקבות נוספים.

חברות וחברים יקרים, המפה הזו היא פרי מאמץ משותף ולבבות פתוחים. ראו בה סקיצה של הנוף ולא תכנית מדויקת. אם עבודתנו מדברת אליכם, הצטרפו אלינו — יחד נשפר אותה.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'legal_title', 'מידע משפטי') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'legend_attention', 'לתשומת לב') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'legend_button_tooltip', 'פתחו את המקרא לרמות הקרינה.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'legend_danger', 'סכנה') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'legend_full_he', 'סולם זה מראה עד כמה מקום בטוח לחיים, למים ולמזון.
זכור: המדידות עשויות להיות חלקיות; השתמש בזה כהכוונה בלבד.

ירוק (0–11 µR/h)
רקע כמעט טבעי.
• מי באר לרוב בטוחים.
• ניתן לגדל צמחים בלי בדיקות.

צהוב (11–30 µR/h)
רקע מוגבר.
• בדוק מים ואדמה.
• בדוק כל מזון לפני אכילה.

אדום (30–100 µR/h)
זיהום חמור.
• אל תשתה את המים.
• גידול או אכילה כאן מסוכנים; בדיקות מעבדה חובה.

שחור (>100 µR/h)
אזור קריטי.
• מים ומזון אינם שימושיים.
• שהייה ממושכת אסורה; רק ביקור קצר עם הגנה.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'legend_safe', 'בטוח') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'legend_title', 'מקרא') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'license_full', 'הפרויקט הזה פועל תחת רישיון <a href="/LICENSE" target="_blank">MIT License</a>. הנוסח המלא נמצא בשורש המאגר ובאתר MIT. מותר לך ללמוד, לשתף ולשנות את הקוד, כל עוד החירויות האלה ממשיכות עם עבודתך. נתוני המחקר מתפרסמים ברישיון Creative Commons 1.0 כדי שהמדידות יישארו בנחלת הכלל.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'license_title', 'רישיון') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_chart_all', 'כל הקריאות') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_chart_averaged', 'ממוצע על פני [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_chart_close', 'סגור') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_chart_day', '24 השעות האחרונות') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_chart_month', '30 הימים האחרונים') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_chart_link', 'פתח תרשימי קרינה') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_country', 'מדינה') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_desc_generic', 'חיישן Safecast זה') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_desc_location', 'מדווח מ-[[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_desc_location_unknown', 'מדווח מאזור לא ידוע') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_desc_named', 'מכשיר [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_desc_no_tube', 'עוקב אחר רמות קרינה.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_desc_radiation_sensor', 'זהו חיישן קרינה.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_desc_transport_air', 'בטיסה') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_desc_transport_bike', 'באופניים') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_desc_transport_car', 'ברכב') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_desc_transport_unknown', 'נייח') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_desc_transport_walk', 'ברגל') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_desc_tube', 'עם גלאי [[tube]].') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_detector', 'גלאי') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_device_id', 'מזהה מכשיר') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_device_label', 'מכשיר') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_extra_intro', 'סביבה') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_history_error', 'לא ניתן לטעון היסטוריה כעת.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_humidity', 'לחות') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_last_seen', 'קריאה אחרונה') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_local_time', 'שעה מקומית עכשיו') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_no_data', 'לא נרשמו נתונים בתקופה זו.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_pressure', 'לחץ אוויר') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_temperature', 'טמפרטורת אוויר') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_temperature_f', 'טמפרטורת אוויר (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_title', 'חיישן Safecast בזמן אמת') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_transport', 'תנועה') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_transport_air', 'מטוס') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_transport_bike', 'אופניים או קורקינט') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_transport_car', 'רכב או טנדר') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_transport_unknown', 'לא זוהה') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'live_marker_transport_walk', 'ברגל') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'locate_button_tooltip', 'מרכז את המפה על המיקום שלי') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'location_error', 'אירעה שגיאה בעת קבלת המיקום.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'location_permission_denied', 'הגישה למיקום נדחתה.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'location_timeout', 'תם הזמן לבקשת המיקום.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'location_unavailable', 'המיקום אינו זמין.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'no_results_found', 'לא נמצאו תוצאות') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'processing_complete', 'העיבוד הושלם!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'processing_on_server', 'מעבד בשרת...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'qr_button_tooltip', 'קוד QR לקישור לאזור מפה זה.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'radiation_dose', 'קצב מנה') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'search_error', 'שגיאת חיפוש. נסה שוב.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'search_input_tooltip', 'חפש מקום על ידי הקלדת האותיות הראשונות של שמו. תופיע רשימת הצעות.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'select_files', 'אנא בחר לפחות קובץ אחד') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'short_link_tooltip', 'לחצו כדי להעתיק קישור קצר לשיתוף') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'sources_full', 'אנחנו מודים לכל מי שמשתף מדידות.

העלאות אנונימיות מציירות שבילים שקטים על המפה.
<a href="https://safecast.org" target="_blank">Safecast</a> מטפחת ארכיון עולמי של קריאות.
<a href="https://atomfast.net" target="_blank">Atomfast</a> שומרת על Atomcloud פועלת.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> אוספת תובנות מ‑Radiacode.

כל תרומה מרחיבה את התמונה המשותפת; נשמח מאוד שתוסיפו את שלכם.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'sources_title', 'מקורות הנתונים') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'speed', 'מהירות') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'speed_filter_tooltip_accuracy', 'מדידות איטיות נשארות קרובות יותר לקרקע, ולכן נתוני הולכי הרגל מדויקים ביותר.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'speed_filter_tooltip_car', 'רכב: 7–200 קמ״ש לנסיעות ומדידות ניידות.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'speed_filter_tooltip_intro', 'בחרו אילו מדידות יוצגו לפי מהירות התנועה.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'speed_filter_tooltip_live', 'לב Safecast: נתוני זמן אמת מ-safecast.org. השתמשו בתיבה כדי להציג או להסתיר מדידות חיות.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'speed_filter_tooltip_ped', 'הולך רגל: פחות מ־7 קמ״ש למדידות בהליכה או בעמידה.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'speed_filter_tooltip_plane', 'מטוס: 200+ קמ״ש למדידות מהאוויר.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'speed_filter_tooltip_title', 'מסנני מהירות') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'theme_toggle_tooltip', 'החליפו בין ערכת המפה הבהירה לכהה.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'title', 'מפת האיזוטופים של Safecast — מפת קרינה של כוכב הלכת הארץ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'track_id', 'מזהה מסלול') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'upload_button', 'העלה [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'upload_button_tooltip', 'הוסף את מסלול המדידות שלך למפה. פורמטים נתמכים: ‎.kml, .kmz, .gpx, .csv, .rctrk, .json, .log‎. ניתן להעלות מספר קבצים, ולאחר ההעלאה תיפתח עמודת המסלול.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'upload_error', 'שגיאה') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'waiting_for_server', 'ממתין לשרת...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('he', 'your_location', 'המיקום שלך') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_archive_desc', 'JSON संग्रह सक्षम होने पर सभी प्रकाशित .json फ़ाइलों वाला tgz पैकेज डाउनलोड करता है।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_archive_link', 'साप्ताहिक संग्रह डाउनलोड करें') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_archive_note', 'संग्रह बंद होने पर सर्वर HTTP 404 Not Found लौटाता है।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_archive_title', 'साप्ताहिक संग्रह पैकेज') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_latest_desc', 'दिए गए अक्षांश, देशांतर और मीटर में त्रिज्या के निकटतम नवीनतम बिंदु लौटाता है।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_latest_link', 'टोक्यो के पास नवीनतम') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_latest_note', 'अपने क्षेत्र पर ध्यान केंद्रित करने के लिए lat, lon और radius_m समायोजित करें।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_latest_title', 'पास की नवीनतम माप') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_root_desc', 'मेटाडेटा, डेटासेट आँकड़े और अन्य सभी एन्डपॉइंट के लिंक सूचीबद्ध करता है।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_root_link', '/api खोलें') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_root_note', 'संग्रहों और सर्वर की स्थिति जानने के लिए यहाँ से शुरू करें।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_root_title', 'मशीन-पठनीय सूचकांक') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_track_index_desc', 'N-वाँ ट्रैक प्राप्त करता है और /api/track/{id}.json जैसा ही JSON लौटाता है।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_track_index_link', 'इंडेक्स 1 प्राप्त करें') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_track_index_note', 'दूसरी प्रविष्टि के लिए संख्या बदलें। प्रतिक्रियाएँ पंक्ति-दर-पंक्ति JSON के रूप में स्ट्रीम होती हैं।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_track_index_title', 'सूचकांक से ट्रैक डाउनलोड करें') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_tracks_desc', 'प्रकाशित ट्रैकों की सूची नाम, संख्या और डाउनलोड लिंक सहित स्ट्रीम करता है।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_tracks_link', 'पहले ट्रैक सूचीबद्ध करें') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_tracks_note', 'लंबी सूचियों को पृष्ठों में बांटने के लिए limit और offset पैरामीटर का उपयोग करें।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_example_tracks_title', 'सभी ट्रैक सारांश') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_examples_heading', 'विशेष एन्डपॉइंट') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_examples_note', 'सभी लिंक JSON में जवाब देते हैं। जब एन्डपॉइंट डेटा स्ट्रीम करते हैं तो ब्राउज़र बड़े फ़ाइलें डाउनलोड कर सकता है।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_intro', 'JSON API मानचित्र पर दिखाए गए डेटा को प्रतिबिंबित करता है। हर उदाहरण नई टैब में खुलता है ताकि आप कच्ची प्रतिक्रिया देख सकें।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_more_docs', 'अधिक विस्तृत मार्गदर्शन चाहिए?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_more_docs_link_label', 'विस्तृत गाइड खोलें') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'api_title', 'API त्वरित प्रारंभ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'attr_legal', 'कानूनी जानकारी') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'attr_license', 'लाइसेंस') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'attr_sources', 'डेटा स्रोत') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'back_to_all_tracks', 'सभी ट्रैकों के संयुक्त मानचित्र पर वापस जाएँ।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'count_rate', 'गणना दर') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'date_slider_tooltip', 'मापन को तारीख या वर्ष सीमा से फ़िल्टर करें.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'date_time', 'तारीख और समय') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'description', 'चिचा का विकिरण मानचित्र') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'download_track_cim', 'ट्रैक डाउनलोड करें (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'duration_days', '[[count]] दिन') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'duration_hours', '[[count]] घंटे') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'duration_minutes', '[[count]] मिनट') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'duration_months', '[[count]] महीने') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'duration_weeks', '[[count]] सप्ताह') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'error_during_upload', 'अपलोड के दौरान त्रुटि!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'error_processing_files', 'फ़ाइलों को प्रोसेस करने में त्रुटि!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'files_uploaded', 'फ़ाइलें अपलोड हो गईं') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'geolocation_not_supported', 'आपका ब्राउज़र भू-स्थान का समर्थन नहीं करता है।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'github_link_tooltip_desc', 'समुदाय द्वारा संचालित मुक्त स्रोत परियोजना.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'github_link_tooltip_title', 'GitHub पर Safecast New Map') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'github_link_tooltip_version', 'वर्तमान संस्करण: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'home_explore_global', 'छोड़ें और वैश्विक मानचित्र देखें') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'home_or', 'या') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'home_search_placeholder', 'शहर, क्षेत्र या देश दर्ज करें...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'home_subtitle', 'शुरू करने के लिए एक स्थान दर्ज करें।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'home_title', 'Safecast विकिरण मानचित्र') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'home_use_location', 'मेरा स्थान उपयोग करें') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'legal_contact', 'प्रतिक्रिया के लिए हमें यहाँ लिखें:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'legal_full', 'प्रस्तावना. हम एक खुला नक्शा बना रहे हैं जहाँ दुनिया भर के लोग डोज़ीमीटर रीडिंग साझा करते हैं ताकि विज्ञान, पर्यावरण, शिक्षा और सुरक्षा जैसे सामूहिक हित मजबूत हों। जब आप अपने आँकड़े प्रकाशित करते हैं तो आप कई लोगों की मदद करते हैं। कृपया इस सामूहिक प्रयास की देखभाल करें।

1) जिम्मेदारी. आप जो जानकारी भेजते हैं उसकी सटीकता और सामग्री के लिए जिम्मेदारी आपकी है। आँकड़े “जैसे हैं” वैसे ही प्रकाशित और उपयोग किए जाते हैं। सेवा उनकी शुद्धता, पूर्णता या किसी उद्देश्य के लिए उपयुक्तता की पुष्टि या गारंटी नहीं देती और उनके उपयोग से उत्पन्न परिणामों के लिए उत्तरदायी नहीं है।

2) खुलापन और लाइसेंस. जब आप माप, तिथियाँ, अनुमानित स्थान, उपकरण का मॉडल या अन्य विवरण साझा करते हैं तो आप समझते हैं कि वे सभी के लिए उपलब्ध हो जाएंगे और CC0 1.0 (पब्लिक डोमेन) डेटा लाइसेंस के तहत स्वतंत्र रूप से इस्तेमाल किए जा सकते हैं। कोड MIT लाइसेंस के अंतर्गत उपलब्ध रहता है। लेखकीय अधिकार आपके पास रहते हैं; कोई पारिश्रमिक नहीं दिया जाता; तीसरे पक्ष द्वारा आगे वितरण हमारे नियंत्रण से बाहर है।

3) “जैसा है” और बिना सत्यापन के. सामग्री पूर्व समीक्षा के बिना प्रकाशित होती है। हम उपकरणों के अंशांकन या त्रुटियों की अनुपस्थिति की गारंटी नहीं दे सकते। यह जानकारी अनुसंधान हेतु साझा की जाती है और पेशेवर सलाह नहीं है।

4) गोपनीयता और संयम. सुरक्षा और विश्वास बनाए रखने के लिए हम समय और निर्देशांक को सामान्यीकृत कर सकते हैं तथा तकनीकी मेटाडेटा को हटाया या अनाम किया जा सकता है। हम ऐसे सामग्री को छिपा या मिटा सकते हैं जिसे हमारे उचित आकलन में स्पैम, नकली, गैरकानूनी या सेवा में व्यवधानकारी माना जाता है।

5) कुकीज़. साइट केवल एक अल्पकालिक तकनीकी सत्र कुकी का उपयोग करती है; आपकी यात्रा समाप्त होते ही वह मिट जाती है। हम अन्य कोई निशान संग्रहित नहीं करते।

मित्रो, यह नक्शा साझा मेहनत और खुले दिलों का परिणाम है। इसे परिदृश्य की रेखाचित्र की तरह समझें, न कि एकदम सटीक आराखड़े की तरह। यदि हमारा काम आपसे जुड़ता है, तो साथ आइए — मिलकर हम इसे और बेहतर बनाएँगे।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'legal_title', 'कानूनी जानकारी') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'legend_attention', 'ध्यान') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'legend_button_tooltip', 'विकिरण स्तर की विवरणिका खोलें.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'legend_danger', 'खतरा') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'legend_full_hi', 'यह पैमाना बताता है कि कोई जगह जीवन, पानी और भोजन के लिए कितनी सुरक्षित है।
याद रखें: माप अधूरे हो सकते हैं; इसे केवल मार्गदर्शक मानें।

हरा (0–11 µR/h)
लगभग प्राकृतिक पृष्ठभूमि।
• कुएँ का पानी आम तौर पर सुरक्षित।
• बिना जाँच के खेती संभव।

पीला (11–30 µR/h)
पृष्ठभूमि बढ़ी हुई।
• पानी और मिट्टी जाँचें।
• किसी भी खाद्य पदार्थ को खाने से पहले जाँचें।

लाल (30–100 µR/h)
गंभीर प्रदूषण।
• पानी न पिएँ।
• यहाँ उगाना या खाना जोखिम भरा; प्रयोगशाला जाँच आवश्यक।

काला (>100 µR/h)
अत्यंत खतरनाक।
• पानी और खाना उपयोगी नहीं।
• लंबे समय तक रहना मना; सिर्फ थोड़ी देर सुरक्षा के साथ।
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'legend_safe', 'सुरक्षित') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'legend_title', 'दिशा-सूचक') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'license_full', 'यह प्रोजेक्ट <a href="/LICENSE" target="_blank">MIT License</a> लाइसेंस के तहत बढ़ता है। पूरा पाठ रिपोजिटरी की मूल डायरेक्टरी और MIT की साइट पर है। आप कोड का अध्ययन, साझा और संशोधन कर सकते हैं, बशर्ते ये स्वतंत्रताएँ आपके कार्य के साथ बनी रहें। अनुसंधान डेटा Creative Commons 1.0 लाइसेंस के तहत प्रकाशित होता है ताकि माप सार्वजनिक क्षेत्र में बने रहें।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'license_title', 'लाइसेंस') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_chart_all', 'सभी रीडिंग') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_chart_averaged', '[[window]] का औसत') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_chart_close', 'बंद करें') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_chart_day', 'पिछले 24 घंटे') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_chart_month', 'पिछले 30 दिन') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_chart_link', 'विकिरण चार्ट खोलें') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_country', 'देश') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_desc_generic', 'यह Safecast सेंसर') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_desc_location', '[[place]] से रिपोर्ट करता है') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_desc_location_unknown', 'अज्ञात क्षेत्र से रिपोर्ट करता है') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_desc_named', 'उपकरण [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_desc_no_tube', 'विकिरण स्तर की निगरानी करता है।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_desc_radiation_sensor', 'यह एक विकिरण सेंसर है।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_desc_transport_air', 'उड़ान के दौरान') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_desc_transport_bike', 'साइकिल से') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_desc_transport_car', 'कार से') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_desc_transport_unknown', 'स्थिर') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_desc_transport_walk', 'पैदल') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_desc_tube', '[[tube]] डिटेक्टर के साथ।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_detector', 'डिटेक्टर') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_device_id', 'उपकरण ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_device_label', 'उपकरण') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_extra_intro', 'पर्यावरण') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_history_error', 'इस समय इतिहास लोड नहीं हो पा रहा।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_humidity', 'आर्द्रता') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_last_seen', 'नवीनतम रीडिंग') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_local_time', 'स्थानीय समय अभी') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_no_data', 'इस अवधि में कोई डेटा दर्ज नहीं है।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_pressure', 'वायुमंडलीय दबाव') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_temperature', 'हवा का तापमान') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_temperature_f', 'हवा का तापमान (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_title', 'Safecast रीयलटाइम सेंसर') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_transport', 'गति') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_transport_air', 'विमान') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_transport_bike', 'साइकिल या स्कूटर') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_transport_car', 'कार या वैन') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_transport_unknown', 'पता नहीं चला') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'live_marker_transport_walk', 'पैदल') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'locate_button_tooltip', 'मेरे स्थान पर मानचित्र केंद्रित करें') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'location_error', 'स्थान प्राप्त करते समय त्रुटि हुई।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'location_permission_denied', 'स्थान पहुँच अस्वीकृत।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'location_timeout', 'स्थान अनुरोध का समय समाप्त हुआ।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'location_unavailable', 'स्थान उपलब्ध नहीं है।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'no_results_found', 'कोई परिणाम नहीं मिला') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'processing_complete', 'प्रोसेसिंग पूर्ण!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'processing_on_server', 'सर्वर पर प्रोसेसिंग हो रही है...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'qr_button_tooltip', 'इस नक्शे के इस हिस्से के लिंक का QR कोड।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'radiation_dose', 'खुराक दर') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'search_error', 'खोज में त्रुटि। कृपया पुनः प्रयास करें।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'search_input_tooltip', 'नाम के पहले अक्षर टाइप करके स्थान खोजें। सुझावों की सूची दिखाई देगी।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'select_files', 'कृपया कम से कम एक फ़ाइल चुनें') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'short_link_tooltip', 'शेयर करने के लिए छोटा लिंक कॉपी करने हेतु क्लिक करें') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'sources_full', 'हम उन सभी का धन्यवाद करते हैं जो माप साझा करते हैं।

गुमनाम अपलोड नक्शे पर शांत पथ बनाते हैं।
<a href="https://safecast.org" target="_blank">Safecast</a> रीडिंग्स का वैश्विक अभिलेख संभालता है।
<a href="https://atomfast.net" target="_blank">Atomfast</a> Atomcloud को चालू रखता है।
<a href="https://radiaverse.com" target="_blank">Radioverse</a> Radiacode से मिली समझ को जोड़ता है।

हर योगदान साझा तस्वीर को विस्तृत करता है; हम आपको भी सादर आमंत्रित करते हैं कि अपना योगदान जोड़ें।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'sources_title', 'डेटा स्रोत') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'speed', 'गति') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'speed_filter_tooltip_accuracy', 'धीमी माप ज़मीन के सबसे करीब रहती हैं, इसलिए पैदल डेटा सबसे सटीक होता है।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'speed_filter_tooltip_car', 'कार: ड्राइव और मोबाइल मापों के लिए 7–200 किमी/घंटा।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'speed_filter_tooltip_intro', 'यात्रा गति के आधार पर कौन सी माप दिखेंगी, चुनें।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'speed_filter_tooltip_live', 'Safecast हृदय: safecast.org से वास्तविक समय डेटा। इस विकल्प से लाइव माप दिखाएँ या छिपाएँ।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'speed_filter_tooltip_ped', 'पैदल: 7 किमी/घंटा से कम पैदल या स्थिर मापों के लिए।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'speed_filter_tooltip_plane', 'विमान: हवाई सर्वेक्षण के लिए 200+ किमी/घंटा।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'speed_filter_tooltip_title', 'गति फ़िल्टर') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'theme_toggle_tooltip', 'मानचित्र की हल्की और गहरी थीम के बीच बदलें.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'title', 'चिचा समस्थानिक मानचित्र — पृथ्वी ग्रह का विकिरण मानचित्र') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'track_id', 'ट्रैक आईडी') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'upload_button', 'अपलोड करें [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'upload_button_tooltip', 'अपना माप ट्रैक मानचित्र में जोड़ें। समर्थित प्रारूप: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log। आप कई फ़ाइलें अपलोड कर सकते हैं और अपलोड के बाद ट्रैक पेज खुलेगा।') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'upload_error', 'त्रुटि') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'waiting_for_server', 'सर्वर की प्रतीक्षा कर रहे हैं...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hi', 'your_location', 'आपका स्थान') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_archive_desc', 'Letölt egy tgz csomagot az összes közzétett .json fájllal, ha a JSON archívum engedélyezett.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_archive_link', 'Heti archívum letöltése') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_archive_note', 'Ha az archívum ki van kapcsolva, a szerver HTTP 404 Not Found hibát ad.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_archive_title', 'Heti archív csomag') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_latest_desc', 'Visszaadja a megadott szélesség, hosszúság és sugár közelében lévő legújabb pontokat.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_latest_link', 'Legfrissebb Tokió közelében') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_latest_note', 'Állítsa be a lat, lon és radius_m értékeket a saját területéhez.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_latest_title', 'Legfrissebb közeli mérési pontok') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_root_desc', 'Felsorolja a metaadatokat, adatkészlet-statisztikákat és minden további végpont hivatkozását.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_root_link', 'Nyissa meg a /api-t') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_root_note', 'Innen indulva ismerheti meg a gyűjteményeket és a szerver állapotát.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_root_title', 'Géppel olvasható index') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_track_index_desc', 'Lekéri az N-edik útvonalat és ugyanazt a JSON-t adja vissza, mint a /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_track_index_link', 'Index 1 lekérése') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_track_index_note', 'Cserélje ki a számot másik bejegyzéshez. A válasz soronkénti JSON-ként érkezik.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_track_index_title', 'Útvonal letöltése index alapján') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_tracks_desc', 'Publikált útvonalak katalógusát streameli nevekkel, darabszámokkal és letöltési linkekkel.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_tracks_link', 'Első útvonalak listázása') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_tracks_note', 'Használja a limit és offset paramétereket a hosszú listák lapozásához.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_example_tracks_title', 'Minden útvonal összefoglalója') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_examples_heading', 'Kiemelt végpontok') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_examples_note', 'Minden hivatkozás JSON-választ ad. Adatfolyam esetén a böngésző nagy fájlokat tölthet le.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_intro', 'A JSON API a térképen látható adatokat tükrözi. Minden példa új fülön nyílik meg, hogy megnézhesse a nyers választ.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_more_docs', 'Részletesebb leírásra van szüksége?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_more_docs_link_label', 'Részletes útmutató megnyitása') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'api_title', 'API gyorstalpaló') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'attr_legal', 'Jogi információk') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'attr_license', 'Licenc') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'attr_sources', 'Adatforrások') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'back_to_all_tracks', 'Vissza az egyesített útvonaltérképhez.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'count_rate', 'Számlálási sebesség') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'date_slider_tooltip', 'Szűrd a mérési adatokat dátum vagy év szerint.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'date_time', 'Dátum és idő') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'description', 'Safecast sugárzási térképe') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'download_track_cim', 'Pálya letöltése (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'duration_days', '[[count]] nap') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'duration_hours', '[[count]] óra') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'duration_minutes', '[[count]] perc') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'duration_months', '[[count]] hónap') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'duration_weeks', '[[count]] hét') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'error_during_upload', 'Hiba a feltöltés közben!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'error_processing_files', 'Hiba a fájlok feldolgozása során!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'files_uploaded', 'Fájlok feltöltve') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'geolocation_not_supported', 'A böngésződ nem támogatja a geolokációt.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'github_link_tooltip_desc', 'Közösség által gondozott nyílt forráskódú projekt.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'github_link_tooltip_title', 'Safecast New Map a GitHubon') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'github_link_tooltip_version', 'Jelenlegi verzió: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'home_explore_global', 'Ugrás a világtérképre') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'home_or', 'vagy') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'home_search_placeholder', 'Adjon meg egy várost, régiót vagy országot...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'home_subtitle', 'Adjon meg egy helyet a kezdéshez.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'home_title', 'Safecast sugárzási térkép') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'home_use_location', 'Saját helyzet használata') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'legal_contact', 'Ha van visszajelzésed, írd meg ide:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'legal_full', 'Előszó. Nyílt térképet építünk, ahol a világ minden tájáról érkező emberek osztják meg doziméteres méréseiket a közjó érdekében – a tudomány, a környezet, az oktatás és a biztonság támogatására. Amikor közzéteszed az adataidat, sokaknak segítesz. Kérünk, bánj gondosan ezzel a közös munkával.

1) Felelősség. Az általad beküldött információk pontosságáért és tartalmáért te felelsz. Az adatok „ahogy vannak” kerülnek közzétételre és felhasználásra. A szolgáltatás nem igazolja és nem garantálja azok helytállóságát, teljességét vagy bármilyen célra való alkalmasságát, és nem vállal felelősséget a felhasználás következményeiért.

2) Nyitottság és licenc. Amikor méréseket, dátumokat, hozzávetőleges helyeket, eszközmodelleket vagy egyéb részleteket osztasz meg, tudomásul veszed, hogy ezek mindenki számára elérhetővé válnak, és szabadon felhasználhatók a CC0 1.0 (Public Domain) adatlicenc szerint. A kód továbbra is MIT licenc alatt érhető el. A szerzői jog nálad marad; ellenszolgáltatás nincs; a harmadik felek általi további terjesztés nincs a felügyeletünk alatt.

3) „Ahogy van” és ellenőrzés nélkül. A közzétételek előzetes felülvizsgálat nélkül jelennek meg. Nem tudjuk garantálni az eszközök kalibrálását vagy a hibamentességet. Az információkat kutatási célokra osztjuk meg, és nem minősülnek szakmai tanácsadásnak.

4) Adatvédelem és moderálás. A biztonság és a bizalom megőrzése érdekében általánosíthatjuk az időpontokat és koordinátákat, illetve eltávolíthatjuk vagy anonimizálhatjuk a technikai metaadatokat. Elrejthetjük vagy törölhetjük azokat az anyagokat, amelyeket ésszerű megítélésünk szerint spamnek, hamisítványnak, jogsértőnek vagy zavarónak találunk. A jóhiszemű méréseket megbecsüljük és igyekszünk megőrizni.

5) Sütik. Az oldal csak egy rövid élettartamú technikai munkamenet-sütit használ; a látogatás végén eltűnik. Más nyomot nem tárolunk.

Barátaink, ez a térkép közös erőfeszítés és nyitott szívek gyümölcse. Tekintsétek a táj vázlatának, ne tökéletes tervrajznak. Ha megszólít benneteket a munkánk, csatlakozzatok – együtt még jobbá tesszük.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'legal_title', 'Jogi tudnivalók') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'legend_attention', 'Figyelem') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'legend_button_tooltip', 'Nyisd meg a sugárzási szintek jelmagyarázatát.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'legend_danger', 'Veszély') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'legend_full_hu', 'Ez a skála megmutatja, mennyire biztonságos egy hely az életre, vízre és élelemre.
Ne feledd: a mérések hiányosak lehetnek; csak iránymutató.

Zöld (0–11 µR/h)
Szinte természetes háttér.
• A kútvíz többnyire biztonságos.
• Növényeket vizsgálat nélkül lehet termeszteni.

Sárga (11–30 µR/h)
Emelkedett háttér.
• Vizsgáld meg a vizet és a talajt.
• Teszteld az ételt evés előtt.

Piros (30–100 µR/h)
Komoly szennyezés.
• Ne igyál a vízből.
• Itt termelni vagy enni veszélyes; laborvizsgálat kell.

Fekete (>100 µR/h)
Kritikus zóna.
• Víz és élelem használhatatlan.
• Csak rövid tartózkodás védelemmel.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'legend_safe', 'Biztonságos') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'legend_title', 'Jelmagyarázat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'license_full', 'Ez a projekt a <a href="/LICENSE" target="_blank">MIT License</a> licenc alatt működik. A teljes szöveg a repó gyökerében és a MIT oldalán található. Szabadon tanulmányozhatod, megoszthatod és módosíthatod a kódot, feltéve hogy ezek a szabadságok a te munkáddal is együtt maradnak. A kutatási adatokat a Creative Commons 1.0 licenc alatt tesszük közzé, hogy a mérések közkinccsé maradjanak.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'license_title', 'Licenc') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_chart_all', 'Összes mérés') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_chart_averaged', 'Átlag [[window]] alatt') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_chart_close', 'Bezárás') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_chart_day', 'Utolsó 24 óra') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_chart_month', 'Utolsó 30 nap') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_chart_link', 'Sugárzási diagramok megnyitása') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_country', 'Ország') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_desc_generic', 'Ez a Safecast érzékelő') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_desc_location', 'jelent innen: [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_desc_location_unknown', 'ismeretlen területről jelent') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_desc_named', 'Eszköz: [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_desc_no_tube', 'figyeli a sugárzási szinteket.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_desc_radiation_sensor', 'Ez egy sugárzásérzékelő.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_desc_transport_air', 'repülés közben') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_desc_transport_bike', 'kerékpárral') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_desc_transport_car', 'autóval') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_desc_transport_unknown', 'álló helyzetben') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_desc_transport_walk', 'gyalog') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_desc_tube', '[[tube]] detektorral.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_detector', 'Detektor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_device_id', 'Eszközazonosító') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_device_label', 'Eszköz') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_extra_intro', 'Környezet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_history_error', 'Az előzmények jelenleg nem tölthetők be.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_humidity', 'Páratartalom') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_last_seen', 'Legutóbbi mérés') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_local_time', 'Helyi idő most') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_no_data', 'Ebben az időszakban nem rögzítettek adatot.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_pressure', 'Légnyomás') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_temperature', 'Levegő hőmérséklete') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_temperature_f', 'Levegő hőmérséklete (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_title', 'Safecast valós idejű érzékelő') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_transport', 'Mozgás') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_transport_air', 'Repülőgép') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_transport_bike', 'Kerékpár vagy roller') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_transport_car', 'Autó vagy furgon') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_transport_unknown', 'Nem észlelt') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'live_marker_transport_walk', 'Gyalog') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'locate_button_tooltip', 'Központosítsd a térképet a helyzetemre') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'location_error', 'Hiba történt a hely lekérése közben.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'location_permission_denied', 'A helyhozzáférés megtagadva.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'location_timeout', 'A helylekérdezés időtúllépéssel megszakadt.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'location_unavailable', 'A hely nem érhető el.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'no_results_found', 'Nincs találat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'processing_complete', 'A feldolgozás befejeződött!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'processing_on_server', 'Feldolgozás a szerveren...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'qr_button_tooltip', 'QR-kód a hivatkozáshoz ehhez a térképrészhez.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'radiation_dose', 'Dózisteljesítmény') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'search_error', 'Keresési hiba. Kérjük, próbálja újra.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'search_input_tooltip', 'Keressen helyet a név első betűinek beírásával. Javaslatlista jelenik meg.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'select_files', 'Kérjük, válassz legalább egy fájlt') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'short_link_tooltip', 'Kattints a rövid megosztási link másolásához') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'sources_full', 'Köszönet mindenkinek, aki megosztja a mérési adatokat.

Az anonim feltöltések csendes útvonalakat rajzolnak a térképre.
<a href="https://safecast.org" target="_blank">Safecast</a> globális archívumot gondoz a leolvasásokból.
<a href="https://atomfast.net" target="_blank">Atomfast</a> életben tartja az Atomcloudot.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> összegyűjti a Radiacode tapasztalatait.

Minden hozzájárulás szélesíti a közös képet; szeretettel várjuk a te adataidat is.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'sources_title', 'Adatforrások') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'speed', 'Sebesség') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'speed_filter_tooltip_accuracy', 'A lassabb mérések maradnak legközelebb a talajhoz, ezért a gyalogos adatok a legpontosabbak.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'speed_filter_tooltip_car', 'Autó: 7–200 km/h az utakhoz és mobil mérésekhez.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'speed_filter_tooltip_intro', 'Válassza ki, mely mérések jelenjenek meg az utazási sebesség alapján.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'speed_filter_tooltip_live', 'Safecast-szív: valós idejű adatok a safecast.org-ról. A jelölőnégyzet mutatja vagy rejti az élő méréseket.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'speed_filter_tooltip_ped', 'Gyalogos: 7 km/h alatti értékek gyalogos vagy álló mérésekhez.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'speed_filter_tooltip_plane', 'Repülő: 200+ km/h a légi felmérésekhez.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'speed_filter_tooltip_title', 'Sebességszűrők') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'theme_toggle_tooltip', 'Válts a világos és sötét térképtéma között.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'title', 'Safecast izotóptérképe — A Föld radiológiai térképe') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'track_id', 'Útvonalazonosító') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'upload_button', 'Feltöltés [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'upload_button_tooltip', 'Add hozzá a mérési útvonaladat a térképhez. Támogatott formátumok: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Több fájlt is feltölthetsz, és a feltöltés után megnyílik az útvonal oldala.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'upload_error', 'Hiba') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'waiting_for_server', 'Várakozás a szerverre...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('hu', 'your_location', 'A te helyzeted') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_archive_desc', 'Mengunduh paket tgz dengan semua berkas .json yang dipublikasikan saat arsip JSON diaktifkan.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_archive_link', 'Unduh arsip mingguan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_archive_note', 'Jika arsip dinonaktifkan server mengembalikan HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_archive_title', 'Paket arsip mingguan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_latest_desc', 'Mengembalikan penanda terbaru di dekat lintang, bujur, dan jari-jari pencarian yang diberikan dalam meter.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_latest_link', 'Terbaru dekat Tokyo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_latest_note', 'Sesuaikan lat, lon, dan radius_m untuk wilayah Anda.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_latest_title', 'Pengukuran terbaru di sekitar') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_root_desc', 'Mendaftar metadata, statistik dataset, dan tautan ke setiap endpoint lainnya.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_root_link', 'Buka /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_root_note', 'Mulai dari sini untuk menemukan koleksi dan status server.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_root_title', 'Indeks yang dapat dibaca mesin') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_track_index_desc', 'Mengambil rute ke-N dan mengembalikan dokumen JSON yang sama seperti /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_track_index_link', 'Ambil indeks 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_track_index_note', 'Ganti angkanya untuk mengambil entri lain. Respons dialirkan sebagai JSON per baris.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_track_index_title', 'Unduh rute berdasarkan indeks') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_tracks_desc', 'Mengalirkan katalog rute yang dipublikasikan lengkap dengan nama, jumlah, dan tautan unduhan.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_tracks_link', 'Daftar rute pertama') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_tracks_note', 'Gunakan parameter limit dan offset untuk menelusuri daftar panjang.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_example_tracks_title', 'Ringkasan semua rute') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_examples_heading', 'Endpoint unggulan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_examples_note', 'Semua tautan merespons dalam JSON. Saat endpoint mengalirkan data, browser dapat mengunduh berkas berukuran besar.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_intro', 'API JSON mencerminkan data yang ditampilkan di peta. Setiap contoh terbuka di tab baru sehingga Anda bisa meninjau respons mentah.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_more_docs', 'Butuh dokumentasi lebih lengkap?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_more_docs_link_label', 'Buka panduan terperinci') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'api_title', 'Mulai cepat API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'attr_legal', 'Informasi hukum') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'attr_license', 'Lisensi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'attr_sources', 'Sumber data') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'back_to_all_tracks', 'Kembali ke peta lintasan gabungan.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'count_rate', 'Laju hitung') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'date_slider_tooltip', 'Saring pengukuran berdasarkan tanggal atau rentang tahun.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'date_time', 'Tanggal dan waktu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'description', 'Peta radiasi Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'download_track_cim', 'Unduh trek (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'duration_days', '[[count]] hari') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'duration_hours', '[[count]] jam') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'duration_minutes', '[[count]] mnt') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'duration_months', '[[count]] bulan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'duration_weeks', '[[count]] minggu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'error_during_upload', 'Kesalahan saat mengunggah!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'error_processing_files', 'Kesalahan saat memproses berkas!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'files_uploaded', 'Berkas telah diunggah') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'geolocation_not_supported', 'Peramban Anda tidak mendukung geolokasi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'github_link_tooltip_desc', 'Proyek open-source yang dikelola komunitas.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'github_link_tooltip_title', 'Safecast New Map di GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'github_link_tooltip_version', 'Versi saat ini: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'home_explore_global', 'Lewati dan jelajahi peta dunia') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'home_or', 'atau') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'home_search_placeholder', 'Masukkan kota, wilayah, atau negara...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'home_subtitle', 'Masukkan lokasi untuk memulai.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'home_title', 'Peta radiologi Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'home_use_location', 'Gunakan lokasi saya') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'legal_contact', 'Untuk masukan, hubungi:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'legal_full', 'Prakata. Kami membangun peta terbuka tempat orang-orang di seluruh dunia berbagi pembacaan dosimeter demi kebaikan bersama — ilmu pengetahuan, ekologi, pendidikan, dan keselamatan. Dengan memublikasikan data, Anda membantu banyak orang. Tolong jaga karya bersama ini dengan penuh perhatian.

1) Tanggung jawab. Ketepatan dan isi informasi yang Anda kirimkan menjadi tanggung jawab Anda. Data diterbitkan dan digunakan “apa adanya”. Layanan tidak mengonfirmasi atau menjamin kebenaran, kelengkapan, atau kesesuaiannya untuk tujuan apa pun dan tidak bertanggung jawab atas akibat apa pun dari penggunaannya.

2) Keterbukaan dan lisensi. Saat Anda membagikan hasil ukur, tanggal, perkiraan lokasi, model perangkat, atau detail lain, Anda memahami bahwa semuanya akan tersedia bagi semua orang dan dapat digunakan secara bebas di bawah lisensi CC0 1.0 (Public Domain) untuk data. Kode tetap berada di bawah lisensi MIT. Hak cipta tetap milik Anda; tidak ada kompensasi; distribusi lebih lanjut oleh pihak ketiga berada di luar kendali kami.

3) “Apa adanya” dan tanpa verifikasi. Publikasi muncul tanpa peninjauan awal. Kami tidak dapat menjamin kalibrasi instrumen atau bebas dari kesalahan. Informasi dibagikan untuk tujuan riset dan bukan merupakan rekomendasi profesional.

4) Privasi dan moderasi. Demi keamanan dan kepercayaan, peta dapat menggeneralisasi waktu dan koordinat; metadata teknis dapat dihapus atau dianonimkan. Kami dapat menyembunyikan atau menghapus materi yang menurut penilaian wajar kami adalah spam, palsu, melanggar hukum, atau mengganggu layanan. Kami merawat pengukuran yang dilakukan dengan itikad baik dan berupaya menjaganya.

5) Cookie. Situs ini hanya menggunakan cookie sesi teknis berumur pendek; cookie tersebut hilang saat kunjungan berakhir. Kami tidak menyimpan jejak lain.

Sahabat-sahabat, peta ini adalah buah dari upaya bersama dan hati yang terbuka. Anggaplah sebagai sketsa bentang alam, bukan cetak biru yang presisi. Jika pekerjaan kami menyentuh Anda, bergabunglah — bersama kita dapat menjadikannya lebih baik.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'legal_title', 'Informasi hukum') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'legend_attention', 'Perhatian') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'legend_button_tooltip', 'Buka legenda tingkat radiasi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'legend_danger', 'Bahaya') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'legend_full_id', 'Skala ini menunjukkan seberapa aman suatu tempat bagi hidup, air, dan makanan.
Ingat: pengukuran bisa tidak lengkap; pakai sebagai panduan saja.

Hijau (0–11 µR/h)
Latar mendekati alami.
• Air sumur umumnya aman.
• Tanaman bisa ditanam tanpa tes.

Kuning (11–30 µR/h)
Latar meningkat.
• Periksa air dan tanah.
• Uji makanan sebelum dimakan.

Merah (30–100 µR/h)
Kontaminasi serius.
• Jangan minum airnya.
• Menanam atau makan di sini berbahaya; tes laboratorium wajib.

Hitam (>100 µR/h)
Zona kritis.
• Air dan makanan tak dapat digunakan.
• Hanya boleh singgah sebentar dengan perlindungan.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'legend_safe', 'Aman') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'legend_title', 'Legenda') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'license_full', 'Proyek ini berjalan di bawah <a href="/LICENSE" target="_blank">MIT License</a>. Teks lengkapnya ada di direktori akar repositori dan di situs MIT. Anda boleh mempelajari, membagikan, dan mengubah kode selama kebebasan ini tetap menyertai karya Anda. Data riset diterbitkan dengan lisensi Creative Commons 1.0 agar pengukuran tetap berada di ranah publik.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'license_title', 'Lisensi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_chart_all', 'Semua pembacaan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_chart_averaged', 'Rata-rata selama [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_chart_close', 'Tutup') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_chart_day', '24 jam terakhir') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_chart_month', '30 hari terakhir') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_chart_link', 'Buka grafik radiasi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_country', 'Negara') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_desc_generic', 'Sensor Safecast ini') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_desc_location', 'melaporkan dari [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_desc_location_unknown', 'melaporkan dari area yang tidak diketahui') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_desc_named', 'Perangkat [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_desc_no_tube', 'memantau tingkat radiasi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_desc_radiation_sensor', 'Ini adalah sensor radiasi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_desc_transport_air', 'saat terbang') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_desc_transport_bike', 'dengan sepeda') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_desc_transport_car', 'dengan mobil') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_desc_transport_unknown', 'diam di tempat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_desc_transport_walk', 'berjalan kaki') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_desc_tube', 'dengan detektor [[tube]].') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_detector', 'Detektor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_device_id', 'ID Perangkat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_device_label', 'Perangkat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_extra_intro', 'Lingkungan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_history_error', 'Riwayat tidak dapat dimuat saat ini.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_humidity', 'Kelembapan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_last_seen', 'Pembacaan terakhir') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_local_time', 'Waktu lokal sekarang') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_no_data', 'Tidak ada data yang tercatat dalam periode ini.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_pressure', 'Tekanan udara') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_temperature', 'Suhu udara') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_temperature_f', 'Suhu udara (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_title', 'Sensor Safecast waktu nyata') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_transport', 'Pergerakan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_transport_air', 'Pesawat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_transport_bike', 'Sepeda atau skuter') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_transport_car', 'Mobil atau van') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_transport_unknown', 'Tidak terdeteksi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'live_marker_transport_walk', 'Berjalan kaki') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'locate_button_tooltip', 'Pusatkan peta pada lokasiku') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'location_error', 'Terjadi kesalahan saat mendapatkan lokasi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'location_permission_denied', 'Akses lokasi ditolak.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'location_timeout', 'Permintaan lokasi kedaluwarsa.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'location_unavailable', 'Lokasi tidak tersedia.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'no_results_found', 'Tidak ditemukan hasil') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'processing_complete', 'Pemrosesan selesai!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'processing_on_server', 'Memproses di server...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'qr_button_tooltip', 'Kode QR tautan untuk area peta ini.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'radiation_dose', 'Laju dosis') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'search_error', 'Kesalahan pencarian. Silakan coba lagi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'search_input_tooltip', 'Cari tempat dengan mengetik huruf pertama namanya. Daftar saran akan muncul.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'select_files', 'Pilih setidaknya satu berkas') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'short_link_tooltip', 'Klik untuk menyalin tautan pendek berbagi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'sources_full', 'Kami berterima kasih kepada semua orang yang berbagi pengukuran.

Unggahan anonim melukis jalur sunyi di peta.
<a href="https://safecast.org" target="_blank">Safecast</a> merawat arsip pembacaan global.
<a href="https://atomfast.net" target="_blank">Atomfast</a> menjaga Atomcloud tetap menyala.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> menghimpun wawasan dari Radiacode.

Setiap kontribusi memperluas gambaran bersama; kami mengundang Anda untuk menambahkan milik Anda.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'sources_title', 'Sumber data') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'speed', 'Kecepatan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'speed_filter_tooltip_accuracy', 'Pengukuran yang lebih lambat tetap paling dekat dengan tanah, sehingga data pejalan kaki paling akurat.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'speed_filter_tooltip_car', 'Mobil: 7–200 km/jam untuk perjalanan dan pengukuran bergerak.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'speed_filter_tooltip_intro', 'Pilih pengukuran yang ditampilkan berdasarkan kecepatan perjalanan.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'speed_filter_tooltip_live', 'Hati Safecast: data waktu nyata dari safecast.org. Gunakan opsi ini untuk menampilkan atau menyembunyikan bacaan langsung.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'speed_filter_tooltip_ped', 'Pejalan kaki: di bawah 7 km/jam untuk pengukuran saat berjalan atau berhenti.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'speed_filter_tooltip_plane', 'Pesawat: 200+ km/jam untuk survei udara.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'speed_filter_tooltip_title', 'Filter kecepatan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'theme_toggle_tooltip', 'Beralih antara tema peta terang dan gelap.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'title', 'Peta Isotop Safecast — Peta radiologi Planet Bumi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'track_id', 'ID lintasan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'upload_button', 'Unggah [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'upload_button_tooltip', 'Tambahkan lintasan pengukuran Anda ke peta. Format yang didukung: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Anda bisa mengunggah beberapa berkas; setelah unggah, halaman lintasan akan terbuka.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'upload_error', 'Kesalahan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'waiting_for_server', 'Menunggu server...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('id', 'your_location', 'Lokasimu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_archive_desc', 'Scarica un archivio tgz con tutti i file .json pubblicati quando l’archivio JSON è attivo.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_archive_link', 'Scarica archivio settimanale') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_archive_note', 'Se l’archivio è disattivato il server restituisce HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_archive_title', 'Pacchetto archivio settimanale') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_latest_desc', 'Restituisce i punti più recenti vicino a latitudine, longitudine e raggio di ricerca indicati in metri.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_latest_link', 'Recenti vicino Tokyo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_latest_note', 'Regola lat, lon e radius_m per la tua area.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_latest_title', 'Misurazioni recenti nelle vicinanze') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_root_desc', 'Elenca metadati, statistiche del dataset e collegamenti a tutti gli altri endpoint.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_root_link', 'Apri /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_root_note', 'Inizia qui per scoprire le collezioni e lo stato del server.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_root_title', 'Indice leggibile dalla macchina') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_track_index_desc', 'Recupera la traccia N-esima e restituisce lo stesso JSON di /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_track_index_link', 'Recupera indice 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_track_index_note', 'Sostituisci il numero per ottenere una voce diversa. Le risposte vengono trasmesse come JSON riga per riga.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_track_index_title', 'Scarica traccia per indice') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_tracks_desc', 'Trasmette un catalogo delle tracce pubblicate con nomi, conteggi e link per il download.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_tracks_link', 'Elenca le prime tracce') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_tracks_note', 'Usa i parametri limit e offset per sfogliare elenchi lunghi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_example_tracks_title', 'Riepiloghi di tutte le tracce') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_examples_heading', 'Endpoint in evidenza') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_examples_note', 'Tutti i link rispondono in JSON. Quando gli endpoint trasmettono dati, il browser può scaricare file di grandi dimensioni.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_intro', 'L’API JSON rispecchia i dati mostrati sulla mappa. Ogni esempio si apre in una nuova scheda per esaminare la risposta grezza.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_more_docs', 'Serve una documentazione più approfondita?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_more_docs_link_label', 'Apri la guida dettagliata') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'api_title', 'Guida rapida all’API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'attr_legal', 'Informazioni legali') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'attr_license', 'Licenza') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'attr_sources', 'Fonti dati') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'back_to_all_tracks', 'Torna alla mappa combinata delle tracce.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'count_rate', 'Tasso di conteggio') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'date_slider_tooltip', 'Filtra le misurazioni per data o intervallo di anni.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'date_time', 'Data e ora') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'description', 'Mappa della radiazione di Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'download_track_cim', 'Scarica traccia (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'duration_days', '[[count]] giorni') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'duration_hours', '[[count]] ore') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'duration_minutes', '[[count]] min') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'duration_months', '[[count]] mesi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'duration_weeks', '[[count]] settimane') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'error_during_upload', 'Errore durante il caricamento!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'error_processing_files', 'Errore durante l''elaborazione dei file!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'files_uploaded', 'File caricati') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'geolocation_not_supported', 'Il tuo browser non supporta la geolocalizzazione.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'github_link_tooltip_desc', 'Progetto open source mantenuto dalla comunità.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'github_link_tooltip_title', 'Safecast New Map su GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'github_link_tooltip_version', 'Versione attuale: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'home_explore_global', 'Salta ed esplora la mappa mondiale') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'home_or', 'o') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'home_search_placeholder', 'Inserisci una città, regione o paese...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'home_subtitle', 'Inserisci una posizione per iniziare.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'home_title', 'Mappa radiologica Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'home_use_location', 'Usa la mia posizione') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'legal_contact', 'Per feedback scrivici a:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'legal_full', 'Preludio. Stiamo costruendo una mappa aperta in cui persone di tutto il mondo condividono le letture dei dosimetri per il bene comune — scienza, ecologia, educazione e sicurezza. Pubblicando i tuoi dati aiuti molte persone. Ti chiediamo di aver cura di questo lavoro collettivo.

1) Responsabilità. L’esattezza e il contenuto delle informazioni inviate restano sotto la responsabilità dell’autore. I dati vengono pubblicati e utilizzati “così come sono”. Il servizio non conferma né garantisce la loro correttezza, completezza o idoneità a qualsiasi scopo e non si assume responsabilità per gli eventuali risultati del loro uso.

2) Apertura e licenza. Condividendo misurazioni, date, posizione approssimativa, modello dello strumento o altri dettagli, comprendi che tutto diventa accessibile a chiunque e può essere usato liberamente secondo la licenza CC0 1.0 (Public Domain) per i dati. Il codice rimane distribuito con licenza MIT. L’autorialità rimane tua; non è previsto alcun compenso; l’eventuale redistribuzione da parte di terzi non è sotto il nostro controllo.

3) “Così com’è” e senza verifica. Le pubblicazioni compaiono senza revisione preventiva. Non possiamo garantire la calibrazione degli strumenti né l’assenza di errori. Le informazioni sono offerte per finalità di ricerca e non costituiscono una raccomandazione professionale.

4) Privacy e moderazione. Per mantenere sicurezza e fiducia, la mappa può generalizzare orari e coordinate; i metadati tecnici possono essere rimossi o resi anonimi. Possiamo nascondere o eliminare materiali che, secondo il nostro ragionevole giudizio, siano spam, falsificazioni, illeciti o d’intralcio al servizio. Trattiamo con cura le misurazioni in buona fede e ci impegniamo a conservarle.

5) Cookie. Il sito usa solo un cookie tecnico di sessione di breve durata; scompare al termine della visita. Non conserviamo altre tracce.

Amiche e amici, la nostra mappa è frutto di impegno condiviso e di cuori aperti. Consideratela come uno schizzo del territorio, non come un progetto preciso. Se il nostro lavoro ti sta a cuore, unisciti a noi: insieme possiamo renderlo migliore.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'legal_title', 'Informazioni legali') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'legend_attention', 'Attenzione') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'legend_button_tooltip', 'Apri la legenda dei livelli di radiazione.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'legend_danger', 'Pericolo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'legend_full_it', 'Questa scala indica quanto un luogo è sicuro per vita, acqua e cibo.
Ricorda: le misurazioni possono essere incomplete; usala solo come guida.

Verde (0–11 µR/h)
Fondo naturale.
• L’acqua di pozzo è generalmente sicura.
• Si possono coltivare piante senza test.

Giallo (11–30 µR/h)
Fondo elevato.
• Controlla acqua e terreno.
• Verifica ogni alimento prima di mangiare.

Rosso (30–100 µR/h)
Contaminazione seria.
• Non bere l’acqua.
• Coltivare o mangiare qui è rischioso; servono esami di laboratorio.

Nero (>100 µR/h)
Zona critica.
• Acqua e cibo inutilizzabili.
• Permanenza lunga vietata; solo brevi visite con protezione.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'legend_safe', 'Sicuro') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'legend_title', 'Legenda') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'license_full', 'Questo progetto cresce sotto la <a href="/LICENSE" target="_blank">MIT License</a>. Il testo completo si trova nella radice del repository e sul sito MIT. Puoi studiare, condividere e modificare il codice, purché queste libertà accompagnino anche il tuo lavoro. I dati della ricerca sono pubblicati con licenza Creative Commons 1.0, così le misurazioni restano nel dominio pubblico.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'license_title', 'Licenza') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_chart_all', 'Tutte le letture') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_chart_averaged', 'Media su [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_chart_close', 'Chiudi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_chart_day', 'Ultime 24 ore') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_chart_month', 'Ultimi 30 giorni') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_chart_link', 'Apri grafici radiazioni') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_country', 'Paese') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_desc_generic', 'Questo sensore Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_desc_location', 'riporta da [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_desc_location_unknown', 'riporta da un''area sconosciuta') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_desc_named', 'Dispositivo [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_desc_no_tube', 'monitora i livelli di radiazione.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_desc_radiation_sensor', 'Questo è un sensore di radiazioni.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_desc_transport_air', 'in volo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_desc_transport_bike', 'in bicicletta') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_desc_transport_car', 'in auto') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_desc_transport_unknown', 'fermo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_desc_transport_walk', 'a piedi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_desc_tube', 'con rilevatore [[tube]].') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_detector', 'Rilevatore') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_device_id', 'ID dispositivo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_device_label', 'Dispositivo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_extra_intro', 'Ambiente') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_history_error', 'Impossibile caricare lo storico al momento.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_humidity', 'Umidità') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_last_seen', 'Ultima lettura') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_local_time', 'Ora locale attuale') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_no_data', 'Nessun dato registrato in questo periodo.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_pressure', 'Pressione atmosferica') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_temperature', 'Temperatura dell''aria') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_temperature_f', 'Temperatura dell''aria (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_title', 'Sensore Safecast in tempo reale') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_transport', 'Movimento') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_transport_air', 'Aereo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_transport_bike', 'Bicicletta o monopattino') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_transport_car', 'Auto o furgone') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_transport_unknown', 'Non rilevato') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'live_marker_transport_walk', 'A piedi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'locate_button_tooltip', 'Centra la mappa sulla mia posizione') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'location_error', 'Si è verificato un errore durante l''ottenimento della posizione.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'location_permission_denied', 'Accesso alla posizione negato.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'location_timeout', 'Richiesta di posizione scaduta.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'location_unavailable', 'Posizione non disponibile.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'no_results_found', 'Nessun risultato trovato') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'processing_complete', 'Elaborazione completata!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'processing_on_server', 'Elaborazione sul server...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'qr_button_tooltip', 'Codice QR del link a quest''area della mappa.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'radiation_dose', 'Tasso di dose') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'search_error', 'Errore di ricerca. Riprova.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'search_input_tooltip', 'Cerca un luogo digitando le prime lettere del nome. Apparirà un elenco di suggerimenti.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'select_files', 'Seleziona almeno un file') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'short_link_tooltip', 'Clicca per copiare un link breve da condividere') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'sources_full', 'Ringraziamo tutte le persone che condividono le misurazioni.

I caricamenti anonimi tracciano percorsi silenziosi sulla mappa.
<a href="https://safecast.org" target="_blank">Safecast</a> cura un archivio globale di rilevazioni.
<a href="https://atomfast.net" target="_blank">Atomfast</a> mantiene vivo l’Atomcloud.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> raccoglie le intuizioni provenienti da Radiacode.

Ogni contributo amplia l’immagine comune; ti invitiamo con calore ad aggiungere il tuo.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'sources_title', 'Fonti dati') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'speed', 'Velocità') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'speed_filter_tooltip_accuracy', 'Le misurazioni più lente restano più vicine al suolo, quindi i dati a piedi sono i più precisi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'speed_filter_tooltip_car', 'Auto: 7–200 km/h per percorsi e rilevazioni mobili.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'speed_filter_tooltip_intro', 'Scegli quali misurazioni mostrare in base alla velocità di spostamento.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'speed_filter_tooltip_live', 'Cuore Safecast: dati in tempo reale da safecast.org. Usa la casella per mostrare o nascondere le letture live.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'speed_filter_tooltip_ped', 'Pedone: meno di 7 km/h per misure a piedi o da fermi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'speed_filter_tooltip_plane', 'Aereo: 200+ km/h per rilievi aerei.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'speed_filter_tooltip_title', 'Filtri di velocità') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'theme_toggle_tooltip', 'Passa tra il tema chiaro e scuro della mappa.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'title', 'Mappa degli isotopi di Safecast — Mappa radiologica del pianeta Terra') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'track_id', 'ID traccia') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'upload_button', 'Carica [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'upload_button_tooltip', 'Aggiungi il tuo percorso di misurazioni alla mappa. Formati supportati: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Puoi caricare più file e, al termine, si aprirà la pagina del percorso.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'upload_error', 'Errore') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'waiting_for_server', 'In attesa del server...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('it', 'your_location', 'La tua posizione') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_archive_desc', 'JSON アーカイブが有効な場合、公開済みの .json ファイルをすべて含む tgz バンドルをダウンロードします。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_archive_link', '週間アーカイブをダウンロード') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_archive_note', 'アーカイブが無効なときは、サーバーは HTTP 404 Not Found を返します。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_archive_title', '週間アーカイブバンドル') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_latest_desc', '指定した緯度・経度と半径（メートル単位）の周辺にある最新のポイントを返します。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_latest_link', '東京周辺の最新データ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_latest_note', 'lat・lon・radius_m を自分の地域に合わせて調整してください。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_latest_title', '近くの最新測定') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_root_desc', 'メタデータ、データセット統計、その他すべてのエンドポイントへのリンクを一覧表示します。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_root_link', '/api を開く') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_root_note', 'コレクションとサーバーの状態を知るにはここから始めてください。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_root_title', '機械可読インデックス') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_track_index_desc', 'N 番目のトラックを取得し、/api/track/{id}.json と同じ JSON を返します。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_track_index_link', 'インデックス 1 を取得') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_track_index_note', '別のエントリを取得するには番号を変更します。応答は行区切り JSON としてストリーミングされます。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_track_index_title', 'インデックスでトラックをダウンロード') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_tracks_desc', '公開済みトラックのカタログを、名称・件数・ダウンロード URL とともに配信します。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_tracks_link', '最初のトラックを一覧表示') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_tracks_note', '長い一覧をめくるには limit と offset パラメータを使用します。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_example_tracks_title', 'すべてのトラック概要') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_examples_heading', '注目のエンドポイント') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_examples_note', 'すべてのリンクが JSON で応答します。エンドポイントがストリーミングする場合、ブラウザーが大きなファイルをダウンロードすることがあります。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_intro', 'JSON API は地図に表示されるデータをそのまま返します。各サンプルは新しいタブで開き、生データを確認できます。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_more_docs', 'さらに詳しい情報が必要ですか？') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_more_docs_link_label', '詳細ガイドを開く') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'api_title', 'API クイックスタート') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'attr_legal', '法的情報') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'attr_license', 'ライセンス') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'attr_sources', 'データソース') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'back_to_all_tracks', 'すべてのトラックの地図に戻る。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'count_rate', '計数率') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'date_slider_tooltip', '測定値を日付または年の範囲で絞り込みます。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'date_time', '日時') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'description', 'チチャの放射線マップ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'download_track_cim', 'トラックをダウンロード (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'duration_days', '[[count]]日') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'duration_hours', '[[count]]時間') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'duration_minutes', '[[count]]分') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'duration_months', '[[count]]ヶ月') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'duration_weeks', '[[count]]週間') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'error_during_upload', 'アップロード中にエラーが発生しました！') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'error_processing_files', 'ファイル処理中にエラーが発生しました！') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'files_uploaded', 'ファイルをアップロードしました') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'geolocation_not_supported', 'お使いのブラウザはジオロケーションをサポートしていません。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'github_link_tooltip_desc', 'コミュニティが運営するオープンソースプロジェクトです。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'github_link_tooltip_title', 'GitHub の Safecast New Map') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'github_link_tooltip_version', '現在のバージョン: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'home_explore_global', 'スキップして世界地図を探索') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'home_or', 'または') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'home_search_placeholder', '都市、地域、または国を入力...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'home_subtitle', '場所を入力して開始します。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'home_title', 'Safecast放射線マップ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'home_use_location', '現在地を使用') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'legal_contact', 'フィードバックは次のメールへお願いします：') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'legal_full', '前文。私たちは、世界中の人々が放射線量計の測定値を共有するオープンな地図を構築しています。科学、生態学、教育、安全のために。データを公開することで、多くの人々の役に立ちます。この共同の成果を大切にしてください。

1) 責任。提出する情報の正確性と内容は、あなた自身の責任です。データは「現状のまま」公開・使用されます。本サービスは、データの正確性、完全性、または特定の目的への適合性を確認・保証するものではなく、その使用から生じるいかなる結果についても責任を負いません。

2) オープン性とライセンス。線量測定データ、日時、おおよその位置、デバイスモデル、その他の事実情報を共有することにより、それらが誰でも利用可能になり、データは CC0 1.0（パブリックドメイン）ライセンスの下で自由に使用できることを了承するものとします。コードは MIT ライセンスの下で提供されます。著作権はあなたに帰属します。報酬は提供されません。第三者による再配布は当方の管理範囲外です。

3) 「現状のまま」・未検証。公開物は事前審査なしに掲載されます。機器の校正やデータに誤りがないことを保証することはできません。情報は研究目的で共有されるものであり、専門的な推奨を構成するものではありません。

4) プライバシーとモデレーション。地図の安全性と信頼性を保つため、タイムスタンプや座標を一般化し、技術的メタデータを削除または匿名化する場合があります。当方の合理的な判断でスパム、偽造、違法、または破壊的とみなされる資料は、非表示または削除する場合があります。誠実な測定データは大切に取り扱い、保存に努めます。

5) Cookie。本サイトは短期間の技術的セッション Cookie のみを使用しており、訪問終了後に消失します。その他の痕跡は保持しません。

皆さん、この地図は共同の努力とオープンな心の結晶です。地形のスケッチとして捉え、正確な設計図とは見なさないでください。私たちの活動に共感いただけましたら、ぜひご参加ください。一緒により良いものにしていきましょう。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'legal_title', 'リーガル情報') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'legend_attention', '注意') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'legend_button_tooltip', '放射線レベルの凡例を開く。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'legend_danger', '危険') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'legend_full_ja', 'この尺度は、場所が生活・水・食べ物にどれほど安全かを示します。
測定値は不完全な場合があるので、参考程度に。

緑 (0–11 µR/h)
ほぼ自然背景。
• 井戸水は概ね安全。
• 植物は検査なしで栽培可。

黄 (11–30 µR/h)
背景上昇。
• 水と土を確認。
• 食べ物は食前に検査を。

赤 (30–100 µR/h)
深刻な汚染。
• 水は飲めない。
• ここでの栽培・摂取は危険；検査が必須。

黒 (>100 µR/h)
極めて危険。
• 水も食料も使えない。
• 長期滞在は禁止；短時間で防護が必要。
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'legend_safe', '安全') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'legend_title', '凡例') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'license_full', 'このプロジェクトは <a href="/LICENSE" target="_blank">MIT ライセンス</a> の下で提供されています。全文はリポジトリのルートおよび MIT サイトに掲載されています。コードの研究、共有、改変が可能ですが、これらの自由はあなたの成果物にも引き継がれる必要があります。研究データセットは Creative Commons 1.0 (CC0) の下でリリースされており、測定データはパブリックドメインに保持されます。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'license_title', 'ライセンス') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_chart_all', '全測定値') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_chart_averaged', '[[window]]の平均') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_chart_close', '閉じる') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_chart_day', '過去24時間') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_chart_month', '過去30日間') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_chart_link', '放射線グラフを開く') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_country', '国') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_desc_generic', 'このSafecastセンサーは') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_desc_location', '[[place]]からレポート') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_desc_location_unknown', '不明な地域からレポート') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_desc_named', 'デバイス [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_desc_no_tube', '放射線レベルを監視中。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_desc_radiation_sensor', 'これは放射線センサーです。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_desc_transport_air', '飛行中') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_desc_transport_bike', '自転車で移動中') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_desc_transport_car', '車で移動中') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_desc_transport_unknown', '静止中') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_desc_transport_walk', '徒歩で移動中') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_desc_tube', '[[tube]]検出器を使用。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_detector', '検出器') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_device_id', 'デバイスID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_device_label', 'デバイス') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_extra_intro', '環境') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_history_error', '現在、履歴を読み込めません。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_humidity', '湿度') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_last_seen', '最新の測定') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_local_time', '現地時間') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_no_data', 'この期間にデータは記録されていません。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_pressure', '気圧') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_temperature', '気温') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_temperature_f', '気温 (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_title', 'Safecastリアルタイムセンサー') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_transport', '移動') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_transport_air', '航空機') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_transport_bike', '自転車またはスクーター') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_transport_car', '自動車またはバン') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_transport_unknown', '未検出') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'live_marker_transport_walk', '徒歩') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'locate_button_tooltip', '地図を現在地にセンタリング') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'location_error', '位置情報の取得中にエラーが発生しました。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'location_permission_denied', '位置情報へのアクセスが拒否されました。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'location_timeout', '位置情報のリクエストがタイムアウトしました。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'location_unavailable', '位置情報を利用できません。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'no_results_found', '結果が見つかりません') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'processing_complete', '処理完了！') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'processing_on_server', 'サーバーで処理中...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'qr_button_tooltip', 'この地図エリアへのリンクのQRコード。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'radiation_dose', '線量率') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'search_error', '検索エラー。もう一度お試しください。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'search_input_tooltip', '名前の最初の数文字を入力して場所を検索します。候補リストが表示されます。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'select_files', '少なくとも 1 つのファイルを選択してください') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'short_link_tooltip', 'クリックして共有用の短いリンクをコピー') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'sources_full', '測定データを共有してくださるすべての方に感謝いたします。

匿名のアップロードが地図上に静かな道筋を描きます。
<a href="https://safecast.org" target="_blank">Safecast</a> は測定データのグローバルアーカイブを維持しています。
<a href="https://atomfast.net" target="_blank">Atomfast</a> は Atomcloud を稼働させ続けています。
<a href="https://radiaverse.com" target="_blank">Radioverse</a> は Radiacode からの知見を集めています。

一つ一つの貢献が全体像を広げます。ぜひあなたもご参加ください。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'sources_title', 'データ提供元') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'speed', '速度') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'speed_filter_tooltip_accuracy', 'ゆっくりした測定ほど地表に近く、徒歩データが最も精度が高くなります。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'speed_filter_tooltip_car', '車：走行や移動測定用に7～200km/h。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'speed_filter_tooltip_intro', '移動速度ごとに表示する測定値を選択します。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'speed_filter_tooltip_live', 'Safecast ハート：safecast.org からのリアルタイムデータ。チェックを切り替えてライブ測定を表示または非表示にします。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'speed_filter_tooltip_ped', '徒歩：歩行や停止中の測定用に7km/h未満。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'speed_filter_tooltip_plane', '航空機：空からの測定用に200+km/h。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'speed_filter_tooltip_title', '速度フィルター') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'theme_toggle_tooltip', '地図のライト/ダークテーマを切り替えます。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'title', 'チチャの同位体マップ — 地球の放射線マップ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'track_id', 'トラック ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'upload_button', 'アップロード [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'upload_button_tooltip', '測定したトラックを地図に追加します。対応形式：.kml、.kmz、.gpx、.csv、.rctrk、.json、.log。複数のファイルをアップロードでき、完了後にトラックページが開きます。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'upload_error', 'エラー') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'waiting_for_server', 'サーバーを待機中...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ja', 'your_location', 'あなたの位置') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_archive_desc', 'JSON 아카이브가 활성화된 경우 게시된 모든 .json 파일이 포함된 tgz 번들을 다운로드합니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_archive_link', '주간 아카이브 다운로드') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_archive_note', '아카이브가 비활성화되어 있으면 서버가 HTTP 404 Not Found를 반환합니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_archive_title', '주간 아카이브 번들') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_latest_desc', '지정한 위도, 경도 및 반경(미터) 주변의 최신 지점을 반환합니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_latest_link', '도쿄 인근 최신 데이터') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_latest_note', 'lat, lon, radius_m 값을 조정해 내 지역을 살펴보세요.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_latest_title', '가까운 최신 측정값') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_root_desc', '메타데이터, 데이터 세트 통계 및 다른 모든 엔드포인트에 대한 링크를 나열합니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_root_link', '/api 열기') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_root_note', '컬렉션과 서버 상태를 살펴보려면 여기서 시작하세요.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_root_title', '기계 판독용 색인') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_track_index_desc', 'N번째 트랙을 찾아 /api/track/{id}.json 과 동일한 JSON을 반환합니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_track_index_link', '인덱스 1 가져오기') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_track_index_note', '다른 항목을 받으려면 숫자를 바꾸세요. 응답은 줄 단위 JSON으로 스트리밍됩니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_track_index_title', '인덱스로 트랙 다운로드') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_tracks_desc', '공개된 트랙의 이름, 개수, 다운로드 링크를 포함한 카탈로그를 스트리밍합니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_tracks_link', '첫 트랙 나열') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_tracks_note', '긴 목록을 탐색하려면 limit와 offset 매개변수를 사용하세요.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_example_tracks_title', '모든 트랙 요약') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_examples_heading', '추천 엔드포인트') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_examples_note', '모든 링크는 JSON으로 응답합니다. 엔드포인트가 데이터를 스트리밍하면 브라우저가 큰 파일을 다운로드할 수 있습니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_intro', 'JSON API는 지도에 표시되는 데이터를 그대로 제공합니다. 각 예시는 새 탭에서 열려 원본 응답을 확인할 수 있습니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_more_docs', '더 자세한 문서가 필요하신가요?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_more_docs_link_label', '상세 가이드 열기') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'api_title', 'API 빠른 시작') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'attr_legal', '법적 정보') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'attr_license', '라이선스') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'attr_sources', '데이터 출처') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'back_to_all_tracks', '모든 트랙 지도로 돌아가기.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'count_rate', '계수율') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'date_slider_tooltip', '날짜 또는 연도 범위로 측정값을 필터링합니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'date_time', '날짜 및 시간') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'description', '치차의 방사선 지도') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'download_track_cim', '트랙 다운로드 (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'duration_days', '[[count]]일') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'duration_hours', '[[count]]시간') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'duration_minutes', '[[count]]분') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'duration_months', '[[count]]개월') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'duration_weeks', '[[count]]주') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'error_during_upload', '업로드 중 오류!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'error_processing_files', '파일 처리 중 오류!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'files_uploaded', '업로드된 파일') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'geolocation_not_supported', '브라우저에서 지리적 위치를 지원하지 않습니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'github_link_tooltip_desc', '커뮤니티가 유지하는 오픈 소스 프로젝트입니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'github_link_tooltip_title', 'GitHub의 Safecast New Map') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'github_link_tooltip_version', '현재 버전: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'home_explore_global', '건너뛰고 세계 지도 탐색') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'home_or', '또는') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'home_search_placeholder', '도시, 지역 또는 국가를 입력...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'home_subtitle', '시작하려면 위치를 입력하세요.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'home_title', 'Safecast 방사선 지도') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'home_use_location', '내 위치 사용') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'legal_contact', '문의는 다음 이메일로 보내주세요:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'legal_full', '서문. 우리는 전 세계 사람들이 방사선 측정값을 공유하는 오픈 지도를 만들고 있습니다. 과학, 생태학, 교육, 안전을 위해서입니다. 데이터를 공개함으로써 많은 사람들을 도울 수 있습니다. 이 공동의 성과를 소중히 다루어 주세요.

1) 책임. 제출하는 정보의 정확성과 내용에 대한 책임은 본인에게 있습니다. 데이터는 "있는 그대로" 게시되고 사용됩니다. 본 서비스는 데이터의 정확성, 완전성 또는 특정 목적에 대한 적합성을 확인하거나 보장하지 않으며, 그 사용으로 인한 어떠한 결과에 대해서도 책임을 지지 않습니다.

2) 개방성과 라이선스. 방사선 측정 데이터, 날짜, 대략적인 위치, 장치 모델 또는 기타 사실 정보를 공유함으로써, 해당 정보가 모든 사람에게 공개되고 데이터는 CC0 1.0(퍼블릭 도메인) 라이선스 하에 자유롭게 사용될 수 있음을 이해합니다. 코드는 MIT 라이선스 하에 제공됩니다. 저작권은 본인에게 귀속됩니다. 보상은 제공되지 않으며, 제3자에 의한 추가 배포는 당사의 통제 범위를 벗어납니다.

3) "있는 그대로" 및 미검증. 게시물은 사전 검토 없이 공개됩니다. 기기 교정의 정확성이나 오류의 부재를 보장할 수 없습니다. 정보는 연구 목적으로 공유되며 전문적인 권고를 구성하지 않습니다.

4) 개인정보 및 관리. 지도의 안전성과 신뢰성을 유지하기 위해 타임스탬프와 좌표를 일반화하고 기술적 메타데이터를 제거하거나 익명화할 수 있습니다. 합리적인 판단에 따라 스팸, 위조, 불법 또는 파괴적이라고 판단되는 자료는 숨기거나 삭제할 수 있습니다. 진실된 측정 데이터는 소중히 취급하며 보존에 힘씁니다.

5) 쿠키. 본 사이트는 단기 기술 세션 쿠키만 사용하며, 방문이 끝나면 사라집니다. 다른 흔적은 남기지 않습니다.

여러분, 이 지도는 공동의 노력과 열린 마음의 결실입니다. 이것을 지형의 스케치로 생각해 주시고, 정확한 청사진으로 여기지 마세요. 우리의 작업에 공감하신다면 함께해 주세요. 함께 더 나은 것을 만들어 갈 수 있습니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'legal_title', '법적 안내') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'legend_attention', '주의') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'legend_button_tooltip', '방사선 수준 범례를 엽니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'legend_danger', '위험') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'legend_full_ko', '이 척도는 장소가 사람, 물, 음식에 얼마나 안전한지 보여줍니다.
측정값은 불완전할 수 있으니 참고용으로만 쓰세요.

초록 (0–11 µR/h)
자연에 가까운 배경.
• 우물물은 대체로 안전.
• 검사 없이 재배 가능.

노랑 (11–30 µR/h)
배경 상승.
• 물과 토양을 확인.
• 음식은 먹기 전 검사.

빨강 (30–100 µR/h)
심각한 오염.
• 물을 마시지 마세요.
• 여기서 재배·섭취는 위험; 실험실 검사 필요.

검정 (>100 µR/h)
치명적 구역.
• 물과 음식 사용 불가.
• 보호장비로 짧게만 머무르세요.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'legend_safe', '안전') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'legend_title', '범례') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'license_full', '이 프로젝트는 <a href="/LICENSE" target="_blank">MIT 라이선스</a> 하에 제공됩니다. 전체 텍스트는 저장소 루트와 MIT 사이트에서 확인할 수 있습니다. 코드를 연구, 공유, 수정할 수 있으며, 이러한 자유는 귀하의 작업에도 함께 전달되어야 합니다. 연구 데이터셋은 Creative Commons 1.0 (CC0)으로 공개되어 측정 데이터가 퍼블릭 도메인에 유지됩니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'license_title', '라이선스') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_chart_all', '모든 측정값') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_chart_averaged', '[[window]] 평균') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_chart_close', '닫기') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_chart_day', '최근 24시간') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_chart_month', '최근 30일') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_chart_link', '방사선 차트 열기') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_country', '국가') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_desc_generic', '이 Safecast 센서는') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_desc_location', '[[place]]에서 보고') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_desc_location_unknown', '알 수 없는 지역에서 보고') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_desc_named', '장치 [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_desc_no_tube', '방사선 수준을 감시 중입니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_desc_radiation_sensor', '이것은 방사선 센서입니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_desc_transport_air', '비행 중') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_desc_transport_bike', '자전거로 이동 중') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_desc_transport_car', '자동차로 이동 중') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_desc_transport_unknown', '정지 중') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_desc_transport_walk', '도보로 이동 중') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_desc_tube', '[[tube]] 검출기 사용.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_detector', '검출기') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_device_id', '장치 ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_device_label', '장치') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_extra_intro', '환경') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_history_error', '현재 기록을 불러올 수 없습니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_humidity', '습도') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_last_seen', '최근 측정') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_local_time', '현지 시간') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_no_data', '이 기간에 기록된 데이터가 없습니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_pressure', '기압') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_temperature', '기온') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_temperature_f', '기온 (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_title', 'Safecast 실시간 센서') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_transport', '이동') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_transport_air', '항공기') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_transport_bike', '자전거 또는 스쿠터') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_transport_car', '자동차 또는 밴') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_transport_unknown', '감지되지 않음') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'live_marker_transport_walk', '도보') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'locate_button_tooltip', '지도를 내 위치로 중심 맞추기') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'location_error', '위치 정보를 가져오는 중 오류가 발생했습니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'location_permission_denied', '위치 접근이 거부되었습니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'location_timeout', '위치 요청 시간이 초과되었습니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'location_unavailable', '위치를 사용할 수 없습니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'no_results_found', '결과를 찾을 수 없습니다') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'processing_complete', '처리 완료!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'processing_on_server', '서버에서 처리 중...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'qr_button_tooltip', '이 지도 영역 링크의 QR 코드.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'radiation_dose', '선량률') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'search_error', '검색 오류. 다시 시도해 주세요.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'search_input_tooltip', '이름의 첫 글자를 입력하여 장소를 검색하세요. 제안 목록이 표시됩니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'select_files', '파일을 하나 이상 선택하세요') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'short_link_tooltip', '클릭하여 짧은 공유 링크를 복사하세요') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'sources_full', '측정 데이터를 공유해 주시는 모든 분들께 감사드립니다.

익명 업로드가 지도 위에 조용한 발자취를 남깁니다.
<a href="https://safecast.org" target="_blank">Safecast</a>는 측정 데이터의 글로벌 아카이브를 운영합니다.
<a href="https://atomfast.net" target="_blank">Atomfast</a>는 Atomcloud를 밝히고 있습니다.
<a href="https://radiaverse.com" target="_blank">Radioverse</a>는 Radiacode의 인사이트를 수집합니다.

각각의 기여가 공동의 그림을 넓힙니다. 여러분도 참여해 주세요.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'sources_title', '데이터 출처') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'speed', '속도') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'speed_filter_tooltip_accuracy', '속도가 느릴수록 지면과 가까이 유지되므로 보행 데이터가 가장 정확합니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'speed_filter_tooltip_car', '자동차: 주행 및 이동 측정을 위한 7–200km/h.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'speed_filter_tooltip_intro', '이동 속도에 따라 표시할 측정을 선택하세요.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'speed_filter_tooltip_live', 'Safecast 하트: safecast.org의 실시간 데이터. 이 옵션으로 실시간 측정을 표시하거나 숨깁니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'speed_filter_tooltip_ped', '보행: 보행 또는 정지 측정을 위한 7km/h 미만.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'speed_filter_tooltip_plane', '비행기: 항공 탐사용 200+km/h.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'speed_filter_tooltip_title', '속도 필터') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'theme_toggle_tooltip', '지도 밝은 테마와 어두운 테마를 전환합니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'title', '치차 동위원소 지도 — 지구 방사선 지도') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'track_id', '트랙 ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'upload_button', '업로드 [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'upload_button_tooltip', '측정 트랙을 지도에 추가하세요. 지원 형식: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. 여러 파일을 업로드할 수 있으며, 업로드가 완료되면 트랙 페이지가 열립니다.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'upload_error', '오류') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'waiting_for_server', '서버 대기 중...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ko', 'your_location', '내 위치') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_archive_desc', 'Memuat turun pakej tgz dengan semua fail .json yang diterbitkan apabila arkib JSON diaktifkan.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_archive_link', 'Muat turun arkib mingguan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_archive_note', 'Jika arkib dilumpuhkan, pelayan memulangkan HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_archive_title', 'Pakej arkib mingguan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_latest_desc', 'Memulangkan penanda terbaharu berhampiran latitud, longitud dan jejari carian yang diberikan dalam meter.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_latest_link', 'Terbaru berhampiran Tokyo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_latest_note', 'Laras nilai lat, lon dan radius_m untuk kawasan anda.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_latest_title', 'Bacaan terbaru berdekatan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_root_desc', 'Menyenaraikan metadata, statistik set data dan pautan ke setiap titik akhir lain.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_root_link', 'Buka /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_root_note', 'Mulakan di sini untuk meneroka koleksi dan status pelayan.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_root_title', 'Indeks boleh baca mesin') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_track_index_desc', 'Mendapatkan trek ke-N dan memulangkan dokumen JSON yang sama seperti /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_track_index_link', 'Dapatkan indeks 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_track_index_note', 'Gantikan nombor untuk mendapatkan entri lain. Respons distrim sebagai JSON baris demi baris.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_track_index_title', 'Muat turun trek mengikut indeks') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_tracks_desc', 'Menstrim katalog trek yang diterbitkan termasuk nama, kiraan dan pautan muat turun.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_tracks_link', 'Senaraikan trek pertama') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_tracks_note', 'Gunakan parameter limit dan offset untuk menyemak senarai panjang.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_example_tracks_title', 'Ringkasan semua trek') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_examples_heading', 'Titik akhir pilihan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_examples_note', 'Semua pautan memberi respons dalam JSON. Apabila endpoint menstrim data, pelayar mungkin memuat turun fail besar.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_intro', 'API JSON memaparkan semula data yang ditunjukkan pada peta. Setiap contoh dibuka pada tab baharu supaya anda boleh melihat respons mentah.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_more_docs', 'Perlu dokumentasi yang lebih terperinci?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_more_docs_link_label', 'Buka panduan terperinci') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'api_title', 'Permulaan pantas API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'attr_legal', 'Maklumat undang-undang') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'attr_license', 'Lesen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'attr_sources', 'Sumber data') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'back_to_all_tracks', 'Kembali ke peta jejak gabungan.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'count_rate', 'Kadar kiraan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'date_slider_tooltip', 'Tapis bacaan mengikut tarikh atau julat tahun.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'date_time', 'Tarikh dan masa') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'description', 'Peta radiasi Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'download_track_cim', 'Muat turun trek (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'duration_days', '[[count]] hari') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'duration_hours', '[[count]] jam') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'duration_minutes', '[[count]] min') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'duration_months', '[[count]] bulan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'duration_weeks', '[[count]] minggu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'error_during_upload', 'Ralat semasa muat naik!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'error_processing_files', 'Ralat semasa memproses fail!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'files_uploaded', 'Fail dimuat naik') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'geolocation_not_supported', 'Pelayar anda tidak menyokong geolokasi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'github_link_tooltip_desc', 'Projek sumber terbuka yang diselenggara oleh komuniti.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'github_link_tooltip_title', 'Safecast New Map di GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'github_link_tooltip_version', 'Versi semasa: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'home_explore_global', 'Langkau dan terokai peta dunia') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'home_or', 'atau') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'home_search_placeholder', 'Masukkan bandar, wilayah, atau negara...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'home_subtitle', 'Masukkan lokasi untuk bermula.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'home_title', 'Peta radiologi Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'home_use_location', 'Gunakan lokasi saya') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'legal_contact', 'Ada maklum balas? Hubungi:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'legal_full', 'Prakata. Kami sedang membangunkan peta terbuka di mana orang ramai dari seluruh dunia berkongsi bacaan dosimeter demi kebaikan bersama — untuk sains, alam sekitar, pendidikan dan keselamatan. Dengan menyiarkan data anda, anda membantu ramai orang. Tolonglah jaga usaha bersama ini dengan penuh rasa tanggungjawab.

1) Tanggungjawab. Ketepatan dan kandungan maklumat yang anda hantar kekal di bawah tanggungjawab anda. Data diterbitkan dan digunakan “seadanya”. Perkhidmatan ini tidak mengesahkan atau menjamin ketepatan, kelengkapan atau kesesuaiannya untuk sebarang tujuan dan tidak menanggung akibat daripada penggunaannya.

2) Keterbukaan dan lesen. Apabila anda berkongsi bacaan, tarikh, lokasi anggaran, model peranti atau butiran lain, anda memahami bahawa ia akan tersedia kepada semua orang dan boleh digunakan secara bebas di bawah lesen data CC0 1.0 (Domain Awam). Kod kekal diterbitkan di bawah lesen MIT. Hak cipta anda tidak berubah; tiada bayaran diberikan; pengedaran semula oleh pihak ketiga berada di luar kawalan kami.

3) “Seadanya” tanpa pengesahan. Sumbangan diterbitkan tanpa semakan awal. Kami tidak dapat menjamin penentukuran peralatan atau ketiadaan kesilapan. Maklumat dikongsi untuk tujuan penyelidikan dan bukan nasihat profesional.

4) Privasi dan moderasi. Demi keselamatan dan kepercayaan, kami mungkin menjadikan masa dan koordinat kurang tepat serta membuang atau menganonimkan metadata teknikal. Kami boleh menyembunyikan atau memadam bahan yang menurut penilaian wajar kami merupakan spam, palsu, menyalahi undang-undang atau mengganggu perkhidmatan. Ukuran yang dihantar dengan niat baik kami pelihara dengan cermat.

5) Kuki. Laman ini hanya menggunakan kuki sesi teknikal yang berjangka pendek; kuki itu hilang sebaik lawatan anda tamat. Kami tidak menyimpan jejak lain.

Rakan-rakan, peta ini lahir daripada usaha bersama dan hati yang terbuka. Anggaplah ia sebagai lakaran rupa bumi, bukan pelan yang serba tepat. Jika kerja kami menyentuh hati anda, sertailah kami — bersama-sama kita boleh menambah baiknya.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'legal_title', 'Maklumat undang-undang') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'legend_attention', 'Perhatian') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'legend_button_tooltip', 'Buka legenda tahap radiasi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'legend_danger', 'Bahaya') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'legend_full_ms', 'Skala ini menunjukkan tahap keselamatan tempat untuk kehidupan, air dan makanan.
Ingat: bacaan mungkin tidak lengkap; guna sebagai panduan sahaja.

Hijau (0–11 µR/h)
Latar hampir semula jadi.
• Air perigi biasanya selamat.
• Boleh menanam tanpa ujian.

Kuning (11–30 µR/h)
Latar meningkat.
• Periksa air dan tanah.
• Uji makanan sebelum makan.

Merah (30–100 µR/h)
Pencemaran serius.
• Jangan minum air.
• Menanam atau makan di sini berbahaya; perlu ujian makmal.

Hitam (>100 µR/h)
Zon kritikal.
• Air dan makanan tidak boleh digunakan.
• Hanya lawatan singkat dengan perlindungan.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'legend_safe', 'Selamat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'legend_title', 'Legenda') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'license_full', 'Projek ini berkembang di bawah <a href="/LICENSE" target="_blank">MIT License</a>. Teks penuh berada di direktori akar repositori dan di laman MIT. Anda boleh mengkaji, berkongsi dan mengubah kod selagi kebebasan ini kekal bersama hasil kerja anda. Data penyelidikan diterbitkan di bawah lesen Creative Commons 1.0 agar bacaan kekal dalam domain awam.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'license_title', 'Lesen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_chart_all', 'Semua bacaan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_chart_averaged', 'Purata dalam [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_chart_close', 'Tutup') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_chart_day', '24 jam terakhir') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_chart_month', '30 hari terakhir') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_chart_link', 'Buka carta sinaran') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_country', 'Negara') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_desc_generic', 'Sensor Safecast ini') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_desc_location', 'melaporkan dari [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_desc_location_unknown', 'melaporkan dari kawasan tidak diketahui') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_desc_named', 'Peranti [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_desc_no_tube', 'memantau tahap sinaran.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_desc_radiation_sensor', 'Ini adalah sensor sinaran.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_desc_transport_air', 'semasa terbang') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_desc_transport_bike', 'dengan basikal') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_desc_transport_car', 'dengan kereta') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_desc_transport_unknown', 'di tempat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_desc_transport_walk', 'berjalan kaki') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_desc_tube', 'dengan pengesan [[tube]].') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_detector', 'Pengesan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_device_id', 'ID Peranti') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_device_label', 'Peranti') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_extra_intro', 'Persekitaran') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_history_error', 'Sejarah tidak dapat dimuat sekarang.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_humidity', 'Kelembapan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_last_seen', 'Bacaan terkini') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_local_time', 'Waktu tempatan sekarang') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_no_data', 'Tiada data direkodkan dalam tempoh ini.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_pressure', 'Tekanan udara') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_temperature', 'Suhu udara') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_temperature_f', 'Suhu udara (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_title', 'Sensor Safecast masa nyata') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_transport', 'Pergerakan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_transport_air', 'Pesawat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_transport_bike', 'Basikal atau skuter') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_transport_car', 'Kereta atau van') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_transport_unknown', 'Tidak dikesan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'live_marker_transport_walk', 'Berjalan kaki') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'locate_button_tooltip', 'Pusatkan peta pada lokasi saya') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'location_error', 'Ralat berlaku semasa mendapatkan lokasi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'location_permission_denied', 'Akses lokasi dinafikan.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'location_timeout', 'Permintaan lokasi tamat masa.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'location_unavailable', 'Lokasi tidak tersedia.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'no_results_found', 'Tiada hasil ditemui') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'processing_complete', 'Pemprosesan selesai!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'processing_on_server', 'Memproses pada pelayan...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'qr_button_tooltip', 'Kod QR pautan untuk kawasan peta ini.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'radiation_dose', 'Kadar dos') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'search_error', 'Ralat carian. Sila cuba lagi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'search_input_tooltip', 'Cari tempat dengan menaip huruf pertama namanya. Senarai cadangan akan muncul.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'select_files', 'Sila pilih sekurang-kurangnya satu fail') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'short_link_tooltip', 'Klik untuk menyalin pautan ringkas perkongsian') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'sources_full', 'Kami berterima kasih kepada semua yang berkongsi bacaan.

Muat naik tanpa nama melakar jejak sunyi di atas peta.
<a href="https://safecast.org" target="_blank">Safecast</a> memelihara arkib global bacaan.
<a href="https://atomfast.net" target="_blank">Atomfast</a> memastikan Atomcloud terus bernyala.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> menghimpun pandangan daripada Radiacode.

Setiap sumbangan meluaskan gambaran bersama; kami mengalu-alukan anda untuk menambah milik anda.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'sources_title', 'Sumber data') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'speed', 'Kelajuan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'speed_filter_tooltip_accuracy', 'Pengukuran yang lebih perlahan kekal paling hampir dengan tanah, jadi data pejalan kaki paling tepat.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'speed_filter_tooltip_car', 'Kereta: 7–200 km/j untuk pemanduan dan tinjauan bergerak.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'speed_filter_tooltip_intro', 'Pilih ukuran yang dipaparkan mengikut kelajuan pergerakan.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'speed_filter_tooltip_live', 'Hati Safecast: data masa nyata dari safecast.org. Gunakan pilihan ini untuk memaparkan atau menyembunyikan bacaan langsung.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'speed_filter_tooltip_ped', 'Pejalan kaki: kurang daripada 7 km/j untuk bacaan ketika berjalan atau berhenti.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'speed_filter_tooltip_plane', 'Pesawat: 200+ km/j untuk tinjauan udara.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'speed_filter_tooltip_title', 'Penapis kelajuan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'theme_toggle_tooltip', 'Tukar antara tema peta cerah dan gelap.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'title', 'Peta Isotop Safecast — Peta radiologi Planet Bumi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'track_id', 'ID jejak') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'upload_button', 'Muat naik [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'upload_button_tooltip', 'Tambahkan jejak pengukuran anda pada peta. Format yang disokong: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Anda boleh memuat naik berbilang fail; selepas muat naik, halaman jejak akan dibuka.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'upload_error', 'Ralat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'waiting_for_server', 'Menunggu pelayan...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ms', 'your_location', 'Lokasi anda') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_archive_desc', 'Downloadt een tgz-pakket met alle gepubliceerde .json-bestanden wanneer het JSON-archief is ingeschakeld.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_archive_link', 'Wekelijks archief downloaden') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_archive_note', 'Als het archief is uitgeschakeld geeft de server HTTP 404 Not Found terug.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_archive_title', 'Wekelijks archiefpakket') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_latest_desc', 'Geeft de nieuwste punten terug in de buurt van de opgegeven breedtegraad, lengtegraad en straal in meters.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_latest_link', 'Nieuwste rond Tokio') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_latest_note', 'Pas lat, lon en radius_m aan om je eigen gebied te bekijken.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_latest_title', 'Nieuwste metingen in de buurt') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_root_desc', 'Toont metadata, datasetstatistieken en links naar alle andere eindpunten.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_root_link', 'Open /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_root_note', 'Begin hier om collecties en de serverstatus te ontdekken.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_root_title', 'Machineleesbare index') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_track_index_desc', 'Haalt de N-de track op en geeft hetzelfde JSON-document als /api/track/{id}.json terug.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_track_index_link', 'Index 1 ophalen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_track_index_note', 'Vervang het getal om een andere vermelding op te halen. Antwoorden worden als regelgescheiden JSON gestreamd.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_track_index_title', 'Track downloaden op index') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_tracks_desc', 'Streamt een catalogus van gepubliceerde tracks met namen, aantallen en downloadlinks.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_tracks_link', 'Eerste tracks tonen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_tracks_note', 'Gebruik de parameters limit en offset om lange lijsten te pagineren.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_example_tracks_title', 'Alle track-samenvattingen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_examples_heading', 'Uitgelichte eindpunten') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_examples_note', 'Alle links geven een JSON-respons. Wanneer eindpunten data streamen kan de browser grote bestanden downloaden.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_intro', 'De JSON-API weerspiegelt de gegevens die op de kaart worden getoond. Elk voorbeeld opent in een nieuw tabblad zodat je de ruwe respons kunt bekijken.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_more_docs', 'Meer documentatie nodig?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_more_docs_link_label', 'Gedetailleerde gids openen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'api_title', 'Snelle start met de API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'attr_legal', 'Juridische informatie') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'attr_license', 'Licentie') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'attr_sources', 'Gegevensbronnen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'back_to_all_tracks', 'Terug naar de gecombineerde sporenkaart.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'count_rate', 'Telsnelheid') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'date_slider_tooltip', 'Filter metingen op datum of jaartal.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'date_time', 'Datum en tijd') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'description', 'Safecast’s stralingskaart') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'download_track_cim', 'Track downloaden (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'duration_days', '[[count]] dagen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'duration_hours', '[[count]] uur') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'duration_minutes', '[[count]] min') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'duration_months', '[[count]] maanden') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'duration_weeks', '[[count]] weken') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'error_during_upload', 'Fout tijdens het uploaden!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'error_processing_files', 'Fout bij het verwerken van de bestanden!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'files_uploaded', 'Bestanden geüpload') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'geolocation_not_supported', 'Uw browser ondersteunt geen geolocatie.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'github_link_tooltip_desc', 'Open-sourceproject onderhouden door de community.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'github_link_tooltip_title', 'Safecast New Map op GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'github_link_tooltip_version', 'Huidige versie: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'home_explore_global', 'Overslaan en de wereldkaart verkennen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'home_or', 'of') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'home_search_placeholder', 'Voer een stad, regio of land in...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'home_subtitle', 'Voer een locatie in om te beginnen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'home_title', 'Safecast stralingkaart') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'home_use_location', 'Gebruik mijn locatie') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'legal_contact', 'Voor feedback, mail naar:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'legal_full', 'Voorwoord. We bouwen aan een open kaart waarop mensen van over de hele wereld hun dosimetermetingen delen voor het algemeen belang — voor wetenschap, milieu, onderwijs en veiligheid. Door je gegevens te publiceren help je velen. We vragen je om zorgvuldig met dit gezamenlijke project om te gaan.

1) Verantwoordelijkheid. Je blijft zelf verantwoordelijk voor de juistheid en inhoud van de informatie die je aanlevert. De gegevens worden “zoals ze zijn” gepubliceerd en gebruikt. De dienst bevestigt of garandeert niet dat ze correct, volledig of geschikt zijn voor welk doel dan ook en aanvaardt geen aansprakelijkheid voor de gevolgen van het gebruik ervan.

2) Openheid en licentie. Door metingen, data, globale locaties, apparaatmodellen of andere details te delen, begrijp je dat deze voor iedereen beschikbaar worden en vrij mogen worden gebruikt onder de CC0 1.0-licentie (Public Domain) voor data. De code blijft beschikbaar onder de MIT-licentie. Het auteurschap blijft bij jou; er staat geen vergoeding tegenover; verdere verspreiding door derden valt buiten onze controle.

3) “Zoals het is” en zonder verificatie. Bijdragen verschijnen zonder voorafgaande beoordeling. We kunnen de kalibratie van apparatuur of de afwezigheid van fouten niet garanderen. De informatie wordt gedeeld voor onderzoeksdoeleinden en vormt geen professioneel advies.

4) Privacy en moderatie. Om veiligheid en vertrouwen te bewaren kunnen we tijden en coördinaten minder precies tonen en technische metadata verwijderen of anonimiseren. We mogen materiaal verbergen of verwijderen dat naar ons redelijk oordeel spam, vervalst, onwettig of storend is. Met metingen die te goeder trouw zijn gedaan gaan we zorgvuldig om en proberen we ze te behouden.

5) Cookies. De site gebruikt alleen een kortlevende technische sessiecookie; deze verdwijnt zodra je bezoek eindigt. We bewaren geen andere sporen.

Beste vrienden, deze kaart is het resultaat van gezamenlijke inspanning en open harten. Zie haar als een schets van het landschap, niet als een perfecte bouwtekening. Voel je je aangesproken door ons werk? Doe mee — samen maken we het nog beter.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'legal_title', 'Juridische informatie') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'legend_attention', 'Let op') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'legend_button_tooltip', 'Open de legenda met stralingsniveaus.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'legend_danger', 'Gevaar') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'legend_full_nl', 'Deze schaal laat zien hoe veilig een plek is voor leven, water en eten.
Onthoud: metingen kunnen onvolledig zijn; gebruik het alleen als richtlijn.

Groen (0–11 µR/h)
Bijna natuurlijke achtergrond.
• Pompwater is meestal veilig.
• Je kunt zonder tests planten kweken.

Geel (11–30 µR/h)
Verhoogde achtergrond.
• Controleer water en bodem.
• Test voedsel voordat je het eet.

Rood (30–100 µR/h)
Ernstige besmetting.
• Drink het water niet.
• Telen of eten van hier is riskant; labtests nodig.

Zwart (>100 µR/h)
Kritieke zone.
• Water en voedsel onbruikbaar.
• Alleen kort verblijf met bescherming.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'legend_safe', 'Veilig') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'legend_title', 'Legenda') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'license_full', 'Dit project valt onder de <a href="/LICENSE" target="_blank">MIT License</a>. De volledige tekst staat in de hoofdmap van de repository en op de MIT-site. Je mag de code bestuderen, delen en aanpassen, zolang deze vrijheden met jouw werk meereizen. De onderzoeksgegevens worden uitgegeven onder de Creative Commons 1.0-licentie, zodat de metingen publiek domein blijven.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'license_title', 'Licentie') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_chart_all', 'Alle metingen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_chart_averaged', 'Gemiddeld over [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_chart_close', 'Sluiten') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_chart_day', 'Laatste 24 uur') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_chart_month', 'Laatste 30 dagen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_chart_link', 'Stralingsgrafieken openen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_country', 'Land') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_desc_generic', 'Deze Safecast-sensor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_desc_location', 'rapporteert vanuit [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_desc_location_unknown', 'rapporteert vanuit een onbekend gebied') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_desc_named', 'Apparaat [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_desc_no_tube', 'bewaakt stralingsniveaus.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_desc_radiation_sensor', 'Dit is een stralingssensor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_desc_transport_air', 'tijdens vliegen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_desc_transport_bike', 'op de fiets') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_desc_transport_car', 'in de auto') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_desc_transport_unknown', 'stilstaand') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_desc_transport_walk', 'te voet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_desc_tube', 'met [[tube]]-detector.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_detector', 'Detector') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_device_id', 'Apparaat-ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_device_label', 'Apparaat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_extra_intro', 'Omgeving') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_history_error', 'Geschiedenis kan momenteel niet worden geladen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_humidity', 'Luchtvochtigheid') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_last_seen', 'Laatste meting') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_local_time', 'Lokale tijd nu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_no_data', 'Geen gegevens geregistreerd in deze periode.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_pressure', 'Luchtdruk') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_temperature', 'Luchttemperatuur') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_temperature_f', 'Luchttemperatuur (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_title', 'Safecast realtime sensor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_transport', 'Beweging') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_transport_air', 'Vliegtuig') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_transport_bike', 'Fiets of step') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_transport_car', 'Auto of bestelbus') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_transport_unknown', 'Niet gedetecteerd') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'live_marker_transport_walk', 'Te voet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'locate_button_tooltip', 'Centreer de kaart op mijn locatie') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'location_error', 'Er is een fout opgetreden bij het ophalen van de locatie.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'location_permission_denied', 'Toegang tot locatie geweigerd.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'location_timeout', 'Time-out voor locatieverzoek.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'location_unavailable', 'Locatie niet beschikbaar.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'no_results_found', 'Geen resultaten gevonden') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'processing_complete', 'Verwerking voltooid!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'processing_on_server', 'Verwerken op de server...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'qr_button_tooltip', 'QR-code van de link naar dit kaartgedeelte.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'radiation_dose', 'Dosisvermogen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'search_error', 'Zoekfout. Probeer het opnieuw.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'search_input_tooltip', 'Zoek een plaats door de eerste letters van de naam te typen. Er verschijnt een suggestielijst.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'select_files', 'Selecteer alstublieft minstens één bestand') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'short_link_tooltip', 'Klik om een korte deellink te kopiëren') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'sources_full', 'Dank aan iedereen die metingen deelt.

Anonieme uploads tekenen stille paden op de kaart.
<a href="https://safecast.org" target="_blank">Safecast</a> onderhoudt een wereldwijd archief van metingen.
<a href="https://atomfast.net" target="_blank">Atomfast</a> houdt de Atomcloud brandend.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> verzamelt inzichten uit Radiacode.

Elke bijdrage verbreedt het gezamenlijke beeld; we nodigen je van harte uit om ook de jouwe toe te voegen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'sources_title', 'Databronnen') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'speed', 'Snelheid') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'speed_filter_tooltip_accuracy', 'Langzamere metingen blijven het dichtst bij de grond; voetgangersgegevens zijn dus het nauwkeurigst.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'speed_filter_tooltip_car', 'Auto: 7–200 km/u voor ritten en mobiele metingen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'speed_filter_tooltip_intro', 'Kies welke metingen worden getoond op basis van de snelheid.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'speed_filter_tooltip_live', 'Safecast-hart: realtime gegevens van safecast.org. Gebruik het selectievakje om live-metingen te tonen of te verbergen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'speed_filter_tooltip_ped', 'Voetganger: onder 7 km/u voor metingen te voet of stilstaand.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'speed_filter_tooltip_plane', 'Vliegtuig: 200+ km/u voor metingen vanuit de lucht.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'speed_filter_tooltip_title', 'Snelheidsfilters') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'theme_toggle_tooltip', 'Schakel tussen het lichte en donkere kaartthema.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'title', 'Safecast Isotopenkaart — Radiologische kaart van planeet Aarde') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'track_id', 'Track-ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'upload_button', 'Uploaden [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'upload_button_tooltip', 'Voeg uw meettraject toe aan de kaart. Ondersteunde formaten: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. U kunt meerdere bestanden uploaden en na het uploaden wordt de trajectpagina geopend.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'upload_error', 'Fout') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'waiting_for_server', 'Wachten op server...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('nl', 'your_location', 'Uw locatie') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_archive_desc', 'Laster ned en tgz-pakke med alle publiserte .json-filer når JSON-arkivet er aktivert.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_archive_link', 'Last ned ukearkiv') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_archive_note', 'Hvis arkivet er deaktivert returnerer serveren HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_archive_title', 'Ukentlig arkivpakke') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_latest_desc', 'Returnerer de nyeste punktene nær angitt breddegrad, lengdegrad og søkeradius i meter.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_latest_link', 'Siste nær Tokyo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_latest_note', 'Juster lat, lon og radius_m for ditt eget område.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_latest_title', 'Siste målinger i nærheten') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_root_desc', 'Lister metadata, datasettstatistikk og lenker til alle andre endepunkter.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_root_link', 'Åpne /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_root_note', 'Start her for å finne samlinger og serverstatus.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_root_title', 'Maskinlesbart register') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_track_index_desc', 'Henter det N-te sporet og returnerer samme JSON som /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_track_index_link', 'Hent indeks 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_track_index_note', 'Bytt ut tallet for å hente en annen oppføring. Svar strømmer som JSON linje for linje.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_track_index_title', 'Last ned spor etter indeks') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_tracks_desc', 'Strømmer en katalog over publiserte spor med navn, antall og nedlastingslenker.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_tracks_link', 'List opp de første sporene') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_tracks_note', 'Bruk parameterne limit og offset for å bla gjennom lange lister.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_example_tracks_title', 'Oversikt over alle spor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_examples_heading', 'Utvalgte endepunkter') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_examples_note', 'Alle lenker svarer med JSON. Når endepunktene strømmer data kan nettleseren laste ned store filer.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_intro', 'JSON-API-en speiler dataene som vises på kartet. Hvert eksempel åpnes i en ny fane slik at du kan inspisere råsvaret.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_more_docs', 'Trenger du mer detaljert dokumentasjon?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_more_docs_link_label', 'Åpne detaljert veiledning') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'api_title', 'Rask start for API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'attr_legal', 'Juridisk informasjon') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'attr_license', 'Lisens') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'attr_sources', 'Datakilder') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'back_to_all_tracks', 'Tilbake til det kombinerte sporkartet.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'count_rate', 'Tellehastighet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'date_slider_tooltip', 'Filtrer målinger etter dato eller år.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'date_time', 'Dato og tid') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'description', 'Safecasts strålingskart') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'download_track_cim', 'Last ned spor (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'duration_days', '[[count]] dager') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'duration_hours', '[[count]] t') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'duration_minutes', '[[count]] min') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'duration_months', '[[count]] måneder') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'duration_weeks', '[[count]] uker') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'error_during_upload', 'Feil under opplasting!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'error_processing_files', 'Feil under behandling av filer!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'files_uploaded', 'Filer lastet opp') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'geolocation_not_supported', 'Nettleseren din støtter ikke geolokalisering.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'github_link_tooltip_desc', 'Et åpent kildekodeprosjekt vedlikeholdt av fellesskapet.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'github_link_tooltip_title', 'Safecast New Map på GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'github_link_tooltip_version', 'Gjeldende versjon: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'home_explore_global', 'Hopp over og utforsk verdenskartet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'home_or', 'eller') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'home_search_placeholder', 'Skriv inn en by, region eller land...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'home_subtitle', 'Skriv inn et sted for å begynne.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'home_title', 'Safecast strålingskart') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'home_use_location', 'Bruk min posisjon') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'legal_contact', 'For tilbakemeldinger, kontakt:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'legal_full', 'Forord. Vi bygger et åpent kart der mennesker fra hele verden deler dosimetermålinger til fellesskapets beste – for vitenskap, miljø, utdanning og sikkerhet. Når du publiserer dataene dine, hjelper du mange. Vi ber deg ta godt vare på dette felles prosjektet.

1) Ansvar. Du er selv ansvarlig for riktigheten og innholdet i informasjonen du sender inn. Dataene publiseres og brukes “som de er”. Tjenesten bekrefter eller garanterer ikke at de er korrekte, fullstendige eller egnet til noe bestemt formål, og den tar ikke ansvar for konsekvenser av bruken.

2) Åpenhet og lisens. Når du deler målinger, datoer, omtrentlige steder, enhetsmodeller eller andre opplysninger, forstår du at de blir tilgjengelige for alle og kan brukes fritt under CC0 1.0-lisensen (Public Domain) for data. Koden forblir tilgjengelig under MIT-lisensen. Opphavsretten forblir hos deg; det gis ingen godtgjørelse; videre distribusjon fra tredjeparter ligger utenfor vår kontroll.

3) “Som de er” og uten forhåndskontroll. Bidrag publiseres uten forhåndsgjennomgang. Vi kan ikke garantere kalibrering av utstyr eller fravær av feil. Informasjonen deles til forskningsformål og utgjør ikke profesjonelle råd.

4) Personvern og moderasjon. For å ivareta sikkerhet og tillit kan vi gjøre tidsangivelser og koordinater mindre presise, og vi kan fjerne eller anonymisere tekniske metadata. Vi kan skjule eller slette materiale som etter vårt rimelige skjønn er spam, forfalsket, ulovlig eller forstyrrende. Ærlige målinger behandler vi med omtanke og forsøker å bevare.

5) Informasjonskapsler. Nettstedet bruker bare en kortlivet teknisk økt-informasjonskapsel; den forsvinner når besøket avsluttes. Vi lagrer ingen andre spor.

Kjære venner, dette kartet er resultatet av felles innsats og åpne hjerter. Se på det som en skisse av terrenget, ikke en millimeternøyaktig plan. Hvis arbeidet vårt gir gjenklang hos deg, bli med — sammen kan vi gjøre det enda bedre.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'legal_title', 'Juridisk informasjon') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'legend_attention', 'Merk!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'legend_button_tooltip', 'Åpne forklaringen for strålingsnivå.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'legend_danger', 'Fare') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'legend_full_no', 'Denne skalaen viser hvor trygg et sted er for liv, vann og mat.
Husk: målingene kan være ufullstendige; bruk den kun som veiledning.

Grønn (0–11 µR/h)
Nær naturlig bakgrunn.
• Brønnvann er vanligvis trygt.
• Planter kan dyrkes uten tester.

Gul (11–30 µR/h)
Bakgrunn økt.
• Sjekk vann og jord.
• Test all mat før du spiser.

Rød (30–100 µR/h)
Alvorlig forurensning.
• Ikke drikk vannet.
• Dyrking eller spising her er farlig; laboratorietester må til.

Svart (>100 µR/h)
Kritisk sone.
• Vann og mat ubrukelige.
• Bare korte besøk med beskyttelse.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'legend_safe', 'Trygt') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'legend_title', 'Forklaring') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'license_full', 'Dette prosjektet lever under <a href="/LICENSE" target="_blank">MIT License</a>. Hele teksten ligger i rotkatalogen til repoet og på MIT-nettstedet. Du kan studere, dele og endre koden så lenge disse frihetene følger arbeidet ditt videre. Forskningsdata publiseres under Creative Commons 1.0 slik at målingene forblir i det offentlige domenet.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'license_title', 'Lisens') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_chart_all', 'Alle målinger') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_chart_averaged', 'Gjennomsnitt over [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_chart_close', 'Lukk') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_chart_day', 'Siste 24 timer') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_chart_month', 'Siste 30 dager') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_chart_link', 'Åpne strålingsdiagrammer') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_country', 'Land') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_desc_generic', 'Denne Safecast-sensoren') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_desc_location', 'rapporterer fra [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_desc_location_unknown', 'rapporterer fra et ukjent område') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_desc_named', 'Enhet [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_desc_no_tube', 'overvåker strålingsnivåer.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_desc_radiation_sensor', 'Dette er en strålingssensor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_desc_transport_air', 'under flyging') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_desc_transport_bike', 'på sykkel') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_desc_transport_car', 'i bil') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_desc_transport_unknown', 'stasjonær') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_desc_transport_walk', 'til fots') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_desc_tube', 'med [[tube]]-detektor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_detector', 'Detektor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_device_id', 'Enhets-ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_device_label', 'Enhet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_extra_intro', 'Miljø') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_history_error', 'Kan ikke laste historikk akkurat nå.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_humidity', 'Luftfuktighet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_last_seen', 'Siste måling') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_local_time', 'Lokal tid nå') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_no_data', 'Ingen data registrert i denne perioden.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_pressure', 'Lufttrykk') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_temperature', 'Lufttemperatur') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_temperature_f', 'Lufttemperatur (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_title', 'Safecast sanntidssensor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_transport', 'Bevegelse') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_transport_air', 'Fly') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_transport_bike', 'Sykkel eller sparkesykkel') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_transport_car', 'Bil eller varebil') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_transport_unknown', 'Ikke oppdaget') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'live_marker_transport_walk', 'Til fots') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'locate_button_tooltip', 'Sentrer kartet på min posisjon') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'location_error', 'Det oppsto en feil ved henting av posisjon.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'location_permission_denied', 'Tilgang til posisjon ble nektet.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'location_timeout', 'Tidsavbrudd for posisjonsforespørsel.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'location_unavailable', 'Posisjon er ikke tilgjengelig.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'no_results_found', 'Ingen resultater funnet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'processing_complete', 'Behandling fullført!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'processing_on_server', 'Behandler på serveren...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'qr_button_tooltip', 'QR-kode for lenken til dette kartområdet.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'radiation_dose', 'Dosehastighet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'search_error', 'Søkefeil. Vennligst prøv igjen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'search_input_tooltip', 'Søk etter et sted ved å skrive de første bokstavene i navnet. En forslagsliste vises.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'select_files', 'Velg minst én fil') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'short_link_tooltip', 'Klikk for å kopiere en kort delingslenke') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'sources_full', 'Vi takker alle som deler målinger.

Anonyme opplastinger tegner stille spor på kartet.
<a href="https://safecast.org" target="_blank">Safecast</a> pleier et verdensomspennende arkiv av måledata.
<a href="https://atomfast.net" target="_blank">Atomfast</a> holder Atomcloud i gang.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> samler innsikt fra Radiacode.

Hver bidrag utvider det felles bildet; vi ønsker ditt bidrag varmt velkommen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'sources_title', 'Datakilder') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'speed', 'Hastighet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'speed_filter_tooltip_accuracy', 'Saktere målinger holder seg nærmest bakken, derfor er data til fots mest presise.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'speed_filter_tooltip_car', 'Bil: 7–200 km/t for kjøreturer og mobile målinger.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'speed_filter_tooltip_intro', 'Velg hvilke målinger som vises etter bevegelseshastighet.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'speed_filter_tooltip_live', 'Safecast-hjerte: sanntidsdata fra safecast.org. Bruk avmerkingsboksen for å vise eller skjule live-målinger.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'speed_filter_tooltip_ped', 'Til fots: under 7 km/t for målinger i gangfart eller i ro.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'speed_filter_tooltip_plane', 'Fly: 200+ km/t for målinger fra luften.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'speed_filter_tooltip_title', 'Fartsfiltre') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'theme_toggle_tooltip', 'Bytt mellom lyst og mørkt karttema.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'title', 'Safecasts isotopkart — Radiologisk kart over planeten Jorden') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'track_id', 'Spor-ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'upload_button', 'Last opp [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'upload_button_tooltip', 'Legg til måleruten din på kartet. Støttede formater: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Du kan laste opp flere filer, og etter opplasting åpnes sporsiden.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'upload_error', 'Feil') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'waiting_for_server', 'Venter på serveren...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('no', 'your_location', 'Din posisjon') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_archive_desc', 'Pobiera paczkę tgz ze wszystkimi opublikowanymi plikami .json, gdy archiwum JSON jest włączone.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_archive_link', 'Pobierz archiwum tygodniowe') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_archive_note', 'Jeśli archiwum jest wyłączone, serwer zwraca błąd HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_archive_title', 'Cotygodniowy pakiet archiwum') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_latest_desc', 'Zwraca najświeższe punkty blisko podanej szerokości, długości i promienia wyszukiwania w metrach.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_latest_link', 'Najnowsze w pobliżu Tokio') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_latest_note', 'Dostosuj lat, lon i radius_m do własnego obszaru.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_latest_title', 'Najnowsze pomiary w pobliżu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_root_desc', 'Wypisuje metadane, statystyki zbioru oraz odsyłacze do wszystkich pozostałych endpointów.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_root_link', 'Otwórz /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_root_note', 'Zacznij tutaj, aby poznać kolekcje i stan serwera.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_root_title', 'Indeks czytelny dla maszyn') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_track_index_desc', 'Pobiera N-tą trasę i zwraca ten sam JSON co /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_track_index_link', 'Pobierz indeks 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_track_index_note', 'Zmień liczbę, aby pobrać inny wpis. Odpowiedzi są strumieniowane jako JSON linia po linii.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_track_index_title', 'Pobierz trasę według indeksu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_tracks_desc', 'Strumieniuje katalog opublikowanych tras z nazwami, liczbą punktów i linkami do pobrania.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_tracks_link', 'Wyświetl pierwsze trasy') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_tracks_note', 'Użyj parametrów limit i offset, aby stronicować długie listy.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_example_tracks_title', 'Podsumowania wszystkich tras') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_examples_heading', 'Polecane endpointy') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_examples_note', 'Wszystkie odnośniki zwracają JSON. Gdy końcówki strumieniują dane, przeglądarka może pobierać duże pliki.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_intro', 'Interfejs JSON odzwierciedla dane pokazywane na mapie. Każdy przykład otwiera się w nowej karcie, aby można było obejrzeć surową odpowiedź.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_more_docs', 'Potrzebujesz dokładniejszej dokumentacji?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_more_docs_link_label', 'Otwórz szczegółowy przewodnik') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'api_title', 'Szybki start API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'attr_legal', 'Informacje prawne') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'attr_license', 'Licencja') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'attr_sources', 'Źródła danych') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'back_to_all_tracks', 'Powrót do łączonej mapy tras.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'count_rate', 'Częstość zliczeń') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'date_slider_tooltip', 'Filtruj pomiary według daty lub zakresu lat.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'date_time', 'Data i czas') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'description', 'Mapa promieniowania Chichy') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'download_track_cim', 'Pobierz trasę (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'duration_days', '[[count]] dni') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'duration_hours', '[[count]] godz.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'duration_minutes', '[[count]] min') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'duration_months', '[[count]] miesięcy') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'duration_weeks', '[[count]] tygodni') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'error_during_upload', 'Błąd podczas przesyłania!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'error_processing_files', 'Błąd podczas przetwarzania plików!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'files_uploaded', 'Pliki przesłane') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'geolocation_not_supported', 'Twoja przeglądarka nie obsługuje geolokalizacji.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'github_link_tooltip_desc', 'Projekt open-source utrzymywany przez społeczność.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'github_link_tooltip_title', 'Safecast New Map na GitHubie') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'github_link_tooltip_version', 'Aktualna wersja: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'home_explore_global', 'Pomiń i przeglądaj mapę świata') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'home_or', 'lub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'home_search_placeholder', 'Wpisz miasto, region lub kraj...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'home_subtitle', 'Wpisz lokalizację, aby rozpocząć.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'home_title', 'Mapa radiologiczna Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'home_use_location', 'Użyj mojej lokalizacji') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'legal_contact', 'W sprawie opinii napisz na:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'legal_full', 'Przedmowa. Tworzymy otwartą mapę, na której ludzie z całego świata dzielą się odczytami dozymetrów dla wspólnego dobra — nauki, środowiska, edukacji i bezpieczeństwa. Publikując swoje dane, pomagasz wielu osobom. Prosimy, traktuj to wspólne dzieło z troską.

1) Odpowiedzialność. Za dokładność i treść przekazywanych informacji odpowiadasz ty. Dane są publikowane i wykorzystywane „tak jak są”. Serwis nie potwierdza ani nie gwarantuje ich poprawności, kompletności ani przydatności do jakiegokolwiek celu i nie ponosi odpowiedzialności za skutki ich użycia.

2) Otwartość i licencja. Udostępniając pomiary, daty, przybliżone lokalizacje, modele urządzeń lub inne szczegóły, przyjmujesz do wiadomości, że staną się dostępne dla wszystkich i mogą być swobodnie wykorzystywane zgodnie z licencją CC0 1.0 (Public Domain) dla danych. Kod pozostaje dostępny na licencji MIT. Autorstwo pozostaje przy tobie; nie wypłacamy wynagrodzenia; dalsza dystrybucja przez osoby trzecie pozostaje poza naszą kontrolą.

3) „Tak jak jest” i bez weryfikacji. Materiały są publikowane bez wcześniejszej kontroli. Nie możemy zagwarantować kalibracji sprzętu ani braku błędów. Informacje udostępniamy w celach badawczych i nie stanowią one porady zawodowej.

4) Prywatność i moderacja. Dla bezpieczeństwa i zaufania możemy uogólniać znaczniki czasu i współrzędne, a także usuwać lub anonimizować techniczne metadane. Możemy ukrywać lub usuwać materiały, które według naszej rozsądnej oceny są spamem, fałszerstwem, naruszają prawo lub zakłócają działanie serwisu. Rzetelne pomiary traktujemy z troską i staramy się je zachować.

5) Pliki cookie. Strona korzysta jedynie z krótkotrwałego technicznego pliku cookie sesji; znika on po zakończeniu wizyty. Nie przechowujemy innych śladów.

Przyjaciele, ta mapa to owoc wspólnego wysiłku i otwartych serc. Traktujcie ją jak szkic terenu, a nie dokładny plan. Jeśli nasza praca do was przemawia, dołączcie — razem uczynimy ją jeszcze lepszą.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'legal_title', 'Informacje prawne') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'legend_attention', 'Uwaga') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'legend_button_tooltip', 'Otwórz legendę poziomów promieniowania.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'legend_danger', 'Zagrożenie') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'legend_full_pl', 'Ta skala pokazuje, na ile miejsce jest bezpieczne dla życia, wody i jedzenia.
Pamiętaj: pomiary mogą być niepełne; traktuj je jedynie orientacyjnie.

Zielony (0–11 µR/h)
Tło prawie naturalne.
• Woda ze studni zazwyczaj bezpieczna.
• Można uprawiać bez testów.

Żółty (11–30 µR/h)
Tło podwyższone.
• Sprawdź wodę i glebę.
• Zbadaj żywność przed spożyciem.

Czerwony (30–100 µR/h)
Poważne skażenie.
• Wody nie pij.
• Uprawa lub jedzenie stąd jest ryzykowne; konieczne badania lab.

Czarny (>100 µR/h)
Strefa krytyczna.
• Woda i jedzenie bezużyteczne.
• Tylko krótki pobyt z ochroną.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'legend_safe', 'Bezpieczne') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'legend_title', 'Legenda') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'license_full', 'Ten projekt rozwija się na licencji <a href="/LICENSE" target="_blank">MIT License</a>. Pełny tekst znajdziesz w katalogu głównym repozytorium oraz na stronie MIT. Możesz badać, udostępniać i modyfikować kod, o ile te wolności pozostają również przy Twojej pracy. Dane badawcze udostępniamy na licencji Creative Commons 1.0, dzięki czemu pomiary pozostają w domenie publicznej.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'license_title', 'Licencja') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_chart_all', 'Wszystkie odczyty') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_chart_averaged', 'Średnia z [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_chart_close', 'Zamknij') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_chart_day', 'Ostatnie 24 godziny') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_chart_month', 'Ostatnie 30 dni') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_chart_link', 'Otwórz wykresy promieniowania') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_country', 'Kraj') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_desc_generic', 'Ten czujnik Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_desc_location', 'raportuje z [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_desc_location_unknown', 'raportuje z nieznanego obszaru') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_desc_named', 'Urządzenie [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_desc_no_tube', 'monitoruje poziom promieniowania.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_desc_radiation_sensor', 'To jest czujnik promieniowania.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_desc_transport_air', 'w locie') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_desc_transport_bike', 'na rowerze') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_desc_transport_car', 'samochodem') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_desc_transport_unknown', 'nieruchomy') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_desc_transport_walk', 'pieszo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_desc_tube', 'z detektorem [[tube]].') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_detector', 'Detektor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_device_id', 'ID urządzenia') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_device_label', 'Urządzenie') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_extra_intro', 'Otoczenie') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_history_error', 'Nie można teraz załadować historii.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_humidity', 'Wilgotność') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_last_seen', 'Ostatni odczyt') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_local_time', 'Czas lokalny teraz') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_no_data', 'Brak danych w tym okresie.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_pressure', 'Ciśnienie atmosferyczne') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_temperature', 'Temperatura powietrza') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_temperature_f', 'Temperatura powietrza (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_title', 'Czujnik Safecast w czasie rzeczywistym') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_transport', 'Ruch') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_transport_air', 'Samolot') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_transport_bike', 'Rower lub hulajnoga') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_transport_car', 'Samochód lub furgonetka') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_transport_unknown', 'Nie wykryto') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'live_marker_transport_walk', 'Pieszo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'locate_button_tooltip', 'Wyśrodkuj mapę na mojej lokalizacji') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'location_error', 'Wystąpił błąd podczas pobierania lokalizacji.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'location_permission_denied', 'Odmowa dostępu do lokalizacji.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'location_timeout', 'Przekroczono czas żądania lokalizacji.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'location_unavailable', 'Lokalizacja niedostępna.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'no_results_found', 'Nie znaleziono wyników') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'processing_complete', 'Przetwarzanie zakończone!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'processing_on_server', 'Przetwarzanie na serwerze...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'qr_button_tooltip', 'Kod QR linku do tego obszaru mapy.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'radiation_dose', 'Moc dawki') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'search_error', 'Błąd wyszukiwania. Spróbuj ponownie.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'search_input_tooltip', 'Szukaj miejsca, wpisując pierwsze litery nazwy. Pojawi się lista sugestii.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'select_files', 'Wybierz co najmniej jeden plik') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'short_link_tooltip', 'Kliknij, aby skopiować krótki link do udostępnienia') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'sources_full', 'Dziękujemy wszystkim, którzy dzielą się pomiarami.

Anonimowe przesłania rysują ciche ścieżki na mapie.
<a href="https://safecast.org" target="_blank">Safecast</a> opiekuje się globalnym archiwum odczytów.
<a href="https://atomfast.net" target="_blank">Atomfast</a> podtrzymuje działanie Atomcloud.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> zbiera wnioski z Radiacode.

Każdy wkład poszerza wspólny obraz; serdecznie zapraszamy, by dodać także swój.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'sources_title', 'Źródła danych') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'speed', 'Prędkość') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'speed_filter_tooltip_accuracy', 'Bardziej powolne pomiary pozostają najbliżej ziemi, dlatego dane piesze są najdokładniejsze.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'speed_filter_tooltip_car', 'Samochód: 7–200 km/h dla przejazdów i mobilnych pomiarów.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'speed_filter_tooltip_intro', 'Wybierz, które pomiary pokazywać według prędkości przemieszczania.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'speed_filter_tooltip_live', 'Serce Safecast: dane w czasie rzeczywistym z safecast.org. Zaznacz/odznacz, aby pokazać lub ukryć pomiary na żywo.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'speed_filter_tooltip_ped', 'Pieszy: poniżej 7 km/h dla pomiarów pieszych lub statycznych.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'speed_filter_tooltip_plane', 'Samolot: 200+ km/h dla pomiarów lotniczych.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'speed_filter_tooltip_title', 'Filtry prędkości') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'theme_toggle_tooltip', 'Przełączaj między jasnym i ciemnym motywem mapy.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'title', 'Mapa izotopów Chichy — Radiologiczna mapa Ziemi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'track_id', 'ID trasy') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'upload_button', 'Prześlij [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'upload_button_tooltip', 'Dodaj swoją trasę pomiarów do mapy. Obsługiwane formaty: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Możesz przesłać wiele plików, a po przesłaniu otworzy się strona trasy.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'upload_error', 'Błąd') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'waiting_for_server', 'Oczekiwanie na serwer...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pl', 'your_location', 'Twoja lokalizacja') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_archive_desc', 'Baixa um pacote tgz com todos os arquivos .json publicados quando o arquivo JSON está ativado.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_archive_link', 'Baixar arquivo semanal') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_archive_note', 'Se o arquivo estiver desativado o servidor retorna HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_archive_title', 'Pacote de arquivo semanal') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_latest_desc', 'Retorna os pontos mais recentes próximos da latitude, longitude e raio pesquisados em metros.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_latest_link', 'Recentes perto de Tóquio') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_latest_note', 'Ajuste lat, lon e radius_m para a sua região.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_latest_title', 'Medições recentes por perto') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_root_desc', 'Lista metadados, estatísticas do conjunto de dados e links para todos os outros endpoints.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_root_link', 'Abrir /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_root_note', 'Comece aqui para descobrir coleções e o status do servidor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_root_title', 'Índice legível por máquina') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_track_index_desc', 'Resolve a N-ésima rota e devolve o mesmo JSON de /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_track_index_link', 'Buscar índice 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_track_index_note', 'Troque o número para obter outra entrada. As respostas são transmitidas como JSON linha a linha.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_track_index_title', 'Baixar rota por índice') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_tracks_desc', 'Transmite um catálogo de rotas publicadas com nomes, contagens e links de download.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_tracks_link', 'Listar primeiras rotas') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_tracks_note', 'Use os parâmetros limit e offset para paginar listas longas.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_example_tracks_title', 'Resumos de todas as rotas') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_examples_heading', 'Endpoints em destaque') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_examples_note', 'Todos os links respondem em JSON. Quando os endpoints transmitem dados, o navegador pode baixar arquivos grandes.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_intro', 'A API JSON reflete os dados exibidos no mapa. Cada exemplo abre em uma nova aba para que você analise a resposta bruta.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_more_docs', 'Precisa de documentação mais completa?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_more_docs_link_label', 'Abrir guia detalhado') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'api_title', 'Início rápido da API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'attr_legal', 'Informações legais') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'attr_license', 'Licença') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'attr_sources', 'Fontes de dados') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'back_to_all_tracks', 'Voltar ao mapa combinado de trilhas.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'count_rate', 'Taxa de contagem') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'date_slider_tooltip', 'Filtre as medições por data ou intervalo de anos.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'date_time', 'Data e hora') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'description', 'Mapa de radiação da Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'download_track_cim', 'Baixar trilha (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'duration_days', '[[count]] dias') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'duration_hours', '[[count]] h') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'duration_minutes', '[[count]] min') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'duration_months', '[[count]] meses') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'duration_weeks', '[[count]] semanas') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'error_during_upload', 'Erro durante o envio!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'error_processing_files', 'Erro ao processar os arquivos!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'files_uploaded', 'Arquivos enviados') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'geolocation_not_supported', 'Seu navegador não suporta geolocalização.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'github_link_tooltip_desc', 'Projeto de código aberto mantido pela comunidade.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'github_link_tooltip_title', 'Safecast New Map no GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'github_link_tooltip_version', 'Versão atual: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'home_explore_global', 'Pular e explorar o mapa mundial') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'home_or', 'ou') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'home_search_placeholder', 'Digite uma cidade, região ou país...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'home_subtitle', 'Digite uma localização para começar.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'home_title', 'Mapa radiológico Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'home_use_location', 'Usar minha localização') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'legal_contact', 'Para feedback, contacte:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'legal_full', 'Preâmbulo. Estamos construindo um mapa aberto no qual pessoas de todo o mundo compartilham leituras de dosímetros pelo bem comum — para a ciência, o meio ambiente, a educação e a segurança. Ao publicar seus dados você ajuda muita gente. Pedimos que cuide com carinho deste trabalho coletivo.

1) Responsabilidade. A precisão e o conteúdo das informações que você envia continuam sob sua responsabilidade. Os dados são publicados e utilizados “como estão”. O serviço não confirma nem garante sua correção, integridade ou adequação para qualquer finalidade e não assume responsabilidade por eventuais consequências do uso.

2) Abertura e licença. Ao compartilhar medições, datas, localizações aproximadas, modelos de dispositivos ou outros detalhes, você entende que eles ficam acessíveis a todos e podem ser utilizados livremente sob a licença de dados CC0 1.0 (Domínio Público). O código permanece disponível sob a licença MIT. A autoria continua sendo sua; não há remuneração; redistribuições futuras por terceiros fogem ao nosso controle.

3) “Como está” e sem verificação prévia. As publicações aparecem sem revisão anterior. Não podemos garantir a calibração dos instrumentos nem a ausência de erros. As informações são compartilhadas para fins de pesquisa e não constituem recomendação profissional.

4) Privacidade e moderação. Para manter a segurança e a confiança, podemos generalizar horários e coordenadas, além de remover ou anonimizar metadados técnicos. Podemos ocultar ou excluir materiais que, a nosso critério razoável, sejam spam, falsificações, ilegais ou perturbem o serviço. Tratamos com cuidado as medições feitas de boa-fé e buscamos preservá-las.

5) Cookies. O site utiliza apenas um cookie técnico de sessão de curta duração; ele desaparece ao final da visita. Não guardamos outros rastros.

Amigas e amigos, este mapa é fruto de um esforço compartilhado e de corações abertos. Encarem-no como um esboço da paisagem, não como um plano milimétrico. Se o nosso trabalho tocar você, venha com a gente — juntos podemos torná-lo ainda melhor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'legal_title', 'Informações legais') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'legend_attention', 'Atenção') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'legend_button_tooltip', 'Abrir a legenda de níveis de radiação.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'legend_danger', 'Perigo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'legend_full_pt', 'Esta escala mostra quão seguro é um lugar para viver, beber e comer.
Lembre: as leituras podem estar incompletas; use apenas como guia.

Verde (0–11 µR/h)
Fundo natural.
• Água de poço geralmente segura.
• Pode-se plantar sem testes.

Amarelo (11–30 µR/h)
Fundo elevado.
• Verifique água e solo.
• Analise qualquer alimento antes de comer.

Vermelho (30–100 µR/h)
Contaminação séria.
• Não beba a água.
• Cultivar ou comer aqui é arriscado; exames laboratoriais obrigatórios.

Preto (>100 µR/h)
Zona crítica.
• Água e comida inutilizáveis.
• Ficar muito tempo é impossível; apenas visitas breves com proteção.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'legend_safe', 'Seguro') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'legend_title', 'Legenda') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'license_full', 'Este projeto cresce sob a <a href="/LICENSE" target="_blank">MIT License</a>. O texto completo está na raiz do repositório e no site da MIT. Você pode estudar, compartilhar e modificar o código, desde que essas liberdades acompanhem o seu trabalho. Os dados de pesquisa são publicados sob a licença Creative Commons 1.0 para que as medições permaneçam em domínio público.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'license_title', 'Licença') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_chart_all', 'Todas as leituras') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_chart_averaged', 'Média em [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_chart_close', 'Fechar') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_chart_day', 'Últimas 24 horas') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_chart_month', 'Últimos 30 dias') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_chart_link', 'Abrir gráficos de radiação') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_country', 'País') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_desc_generic', 'Este sensor Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_desc_location', 'reporta de [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_desc_location_unknown', 'reporta de uma área desconhecida') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_desc_named', 'Dispositivo [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_desc_no_tube', 'monitora os níveis de radiação.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_desc_radiation_sensor', 'Este é um sensor de radiação.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_desc_transport_air', 'em voo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_desc_transport_bike', 'de bicicleta') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_desc_transport_car', 'de carro') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_desc_transport_unknown', 'estacionário') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_desc_transport_walk', 'a pé') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_desc_tube', 'com detector [[tube]].') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_detector', 'Detector') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_device_id', 'ID do dispositivo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_device_label', 'Dispositivo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_extra_intro', 'Ambiente') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_history_error', 'Não foi possível carregar o histórico agora.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_humidity', 'Umidade') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_last_seen', 'Última leitura') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_local_time', 'Hora local atual') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_no_data', 'Nenhum dado registrado neste período.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_pressure', 'Pressão atmosférica') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_temperature', 'Temperatura do ar') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_temperature_f', 'Temperatura do ar (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_title', 'Sensor Safecast em tempo real') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_transport', 'Movimento') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_transport_air', 'Aeronave') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_transport_bike', 'Bicicleta ou patinete') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_transport_car', 'Carro ou van') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_transport_unknown', 'Não detectado') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'live_marker_transport_walk', 'A pé') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'locate_button_tooltip', 'Centralizar o mapa na minha localização') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'location_error', 'Ocorreu um erro ao obter a localização.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'location_permission_denied', 'Acesso à localização negado.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'location_timeout', 'Tempo da solicitação de localização esgotado.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'location_unavailable', 'Localização indisponível.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'no_results_found', 'Nenhum resultado encontrado') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'processing_complete', 'Processamento concluído!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'processing_on_server', 'Processando no servidor...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'qr_button_tooltip', 'Código QR do link para esta área do mapa.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'radiation_dose', 'Taxa de dose') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'search_error', 'Erro na busca. Tente novamente.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'search_input_tooltip', 'Pesquise um lugar digitando as primeiras letras do nome. Uma lista de sugestões aparecerá.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'select_files', 'Selecione pelo menos um arquivo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'short_link_tooltip', 'Clique para copiar um link curto de compartilhamento') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'sources_full', 'Agradecemos a todas as pessoas que compartilham medições.

Envios anônimos traçam caminhos silenciosos no mapa.
<a href="https://safecast.org" target="_blank">Safecast</a> mantém um arquivo global de leituras.
<a href="https://atomfast.net" target="_blank">Atomfast</a> mantém a Atomcloud acesa.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> reúne os aprendizados da Radiacode.

Cada contribuição amplia o panorama comum; convidamos você, de coração, a trazer a sua.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'sources_title', 'Fontes de dados') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'speed', 'Velocidade') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'speed_filter_tooltip_accuracy', 'Medições mais lentas ficam mais próximas do solo; por isso os dados de pedestres são os mais precisos.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'speed_filter_tooltip_car', 'Carro: 7–200 km/h para deslocamentos e medições móveis.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'speed_filter_tooltip_intro', 'Escolha quais medições mostrar de acordo com a velocidade de deslocamento.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'speed_filter_tooltip_live', 'Coração Safecast: dados em tempo real de safecast.org. Use a opção para mostrar ou ocultar as leituras ao vivo.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'speed_filter_tooltip_ped', 'Pedestre: abaixo de 7 km/h para leituras a pé ou parado.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'speed_filter_tooltip_plane', 'Avião: 200+ km/h para medições aéreas.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'speed_filter_tooltip_title', 'Filtros de velocidade') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'theme_toggle_tooltip', 'Alterne entre os temas claro e escuro do mapa.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'title', 'Mapa de Isótopos da Safecast — Mapa radiológico do planeta Terra') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'track_id', 'ID da trilha') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'upload_button', 'Enviar [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'upload_button_tooltip', 'Adicione a sua rota de medições ao mapa. Formatos suportados: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. É possível enviar vários ficheiros e, após o envio, abrirá a página da rota.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'upload_error', 'Erro') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'waiting_for_server', 'Aguardando o servidor...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('pt', 'your_location', 'Sua localização') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_archive_desc', 'Скачивает архив tgz со всеми опубликованными файлами .json, если включён JSON‑архив.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_archive_link', 'Скачать недельный архив') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_archive_note', 'Если архив выключен, сервер вернёт HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_archive_title', 'Еженедельный архив') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_latest_desc', 'Возвращает самые новые точки рядом с указанными широтой, долготой и радиусом поиска в метрах.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_latest_link', 'Свежие рядом с Токио') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_latest_note', 'Настройте lat, lon и radius_m под свой район.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_latest_title', 'Свежие измерения поблизости') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_root_desc', 'Перечисляет метаданные, статистику наборов и ссылки на все остальные эндпоинты.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_root_link', 'Открыть /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_root_note', 'Начните отсюда, чтобы узнать коллекции и статус сервера.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_root_title', 'Машиночитаемый индекс') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_track_index_desc', 'Выдаёт N‑й маршрут и возвращает тот же JSON, что и /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_track_index_link', 'Получить индекс 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_track_index_note', 'Замените число, чтобы получить другую запись. Ответ передаётся построчным JSON.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_track_index_title', 'Скачать маршрут по номеру') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_tracks_desc', 'Стримит каталог опубликованных маршрутов с названиями, количеством и ссылками на скачивание.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_tracks_link', 'Список первых маршрутов') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_tracks_note', 'Используйте параметры limit и offset, чтобы разбивать длинные списки на страницы.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_example_tracks_title', 'Все сводки маршрутов') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_examples_heading', 'Избранные эндпоинты') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_examples_note', 'Все ссылки отвечают в формате JSON. Когда конечные точки отдают поток, браузер может скачивать крупные файлы.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_intro', 'JSON‑API повторяет данные, которые вы видите на карте. Каждый пример открывается в новой вкладке, чтобы можно было посмотреть сырой ответ.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_more_docs', 'Нужны более подробные материалы?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_more_docs_link_label', 'Открыть подробное руководство') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'api_title', 'Быстрый старт API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'attr_legal', 'Правовые сведения') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'attr_license', 'Лицензия') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'attr_sources', 'Источники данных') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'back_to_all_tracks', 'Назад к общей карте треков.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'count_rate', 'Скорость счета') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'date_slider_tooltip', 'Фильтровать измерения по дате или диапазону лет.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'date_time', 'Дата и время') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'description', 'Чичина карта радиации') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'download_track_cim', 'Скачать трек (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'duration_days', '[[count]] дн') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'duration_hours', '[[count]] ч') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'duration_minutes', '[[count]] мин') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'duration_months', '[[count]] мес') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'duration_weeks', '[[count]] нед') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'error_during_upload', 'Ошибка при загрузке!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'error_processing_files', 'Ошибка обработки файлов!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'files_uploaded', 'Файлы загружены') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'geolocation_not_supported', 'Ваш браузер не поддерживает геолокацию.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'github_link_tooltip_desc', 'Проект с открытым исходным кодом, поддерживаемый сообществом.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'github_link_tooltip_title', 'Safecast New Map на GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'github_link_tooltip_version', 'Текущая версия: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'home_explore_global', 'Пропустить и открыть карту мира') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'home_or', 'или') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'home_search_placeholder', 'Введите город, регион или страну...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'home_subtitle', 'Введите место для начала.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'home_title', 'Радиологическая карта Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'home_use_location', 'Использовать моё местоположение') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'legal_contact', 'Для обратной связи напишите:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'legal_full', 'Друзья, вместе с Вами мы создаём открытую карту, где люди со всего мира делятся показаниями дозиметров ради общего дела — науки, экологии, образования и безопасности. Публикуя данные, вы помогаете множеству людей со всего мира. Просим Вас бережно относиться к этому общему труду.

1) Ответственность. Сервис не несет ответственности за публикуемые данные. Точность и содержание отправляемой информации остаются на совести авторов. Данные публикуются и используются «как есть». Сервис не подтверждает и не гарантирует их корректность, полноту или пригодность для какой-либо цели и не несёт ответственности за возможные последствия их использования.

2) Открытость и лицензия. Отправляя измерения, даты, примерные координаты и иные сведения, вы понимаете, что они становятся доступными всем и могут свободно использоваться по лицензии CC0 1.0 (Public Domain) для данных. Код проекта остаётся под лицензией MIT. Авторство сохраняется за вами; вознаграждение не предусмотрено; дальнейшее распространение третьими лицами не контролируется сервисом.

3) «Как есть» и без проверки. Материалы публикуются без предварительной модерации. Мы не можем гарантировать поверку и калибровку приборов, а так же отсутствие ошибок. Информация предназначена для исследовательских целей и не является профессиональной рекомендацией.

4) Конфиденциальность и модерация. Ради безопасности и общего доверия мы можем обобщать время и координаты, удалять или обезличивать технические метаданные. Мы вправе скрывать или удалять материалы, которые по нашему взвешенному мнению являются спамом, подделкой, нарушают закон или мешают работе сервиса. К добросовестным измерениям мы относимся внимательно и стремимся их сохранить.

5) Файлы cookie. Сайт использует только короткоживущий технический cookie сессии; он исчезает после окончания визита. Мы не храним других следов.

Друзья, эта карта — результат совместного труда и открытых сердец. Воспринимайте её как набросок, а не как точный чертёж. Если наша работа вам откликается, присоединяйтесь — вместе мы сделаем её еще лучше.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'legal_title', 'Правовая информация') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'legend_attention', 'Внимание') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'legend_button_tooltip', 'Открыть легенду уровней радиации.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'legend_danger', 'Опасно') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'legend_full_ru', 'Эта шкала показывает, насколько место безопасно для людей, воды и еды.
Помните: измерения могут быть неполными, радиация может быть выше или скрыта. Относитесь к числам как к ориентиру.

Зеленая (0–11 µR/h)
Фон близок к естественному.
• Вода из скважин обычно безопасна.
• Можно выращивать растения без проверок.

Желтая (11–30 µR/h)
Фон повышен; будьте осторожны.
• Проверяйте воду и почву.
• Тестируйте овощи, грибы и другие продукты перед употреблением.

Красная (30–100 µR/h)
Серьёзное загрязнение.
• Воду пить нельзя.
• Выращивать или есть продукты отсюда рискованно; нужны анализы.

Чёрная (>100 µR/h)
Критическая зона.
• Воду и пищу использовать нельзя.
• Долгое пребывание исключено; только короткие визиты с защитой.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'legend_safe', 'Безопасно') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'legend_title', 'Легенда') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'license_full', 'Проект живёт под лицензией <a href="/LICENSE" target="_blank">MIT License</a>. Полный текст лежит в корне репозитория и на сайте MIT. Вы вольны изучать, распространять и менять код, сохраняя эти свободы для других. Исследовательские данные публикуются по лицензии Creative Commons 1.0, чтобы измерения оставались в общественном достоянии.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'license_title', 'Лицензия') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_chart_all', 'Вся история наблюдений') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_chart_averaged', 'Среднее за [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_chart_close', 'Закрыть') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_chart_day', 'Последние 24 часа') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_chart_month', 'Последние 30 дней') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_chart_link', 'Открыть графики радиации') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_country', 'Страна') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_desc_generic', 'Этот датчик Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_desc_location', 'передает данные рядом с [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_desc_location_unknown', 'передает данные из неизвестного места') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_desc_named', 'Прибор [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_desc_no_tube', 'наблюдает уровень радиации.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_desc_radiation_sensor', 'Это датчик радиации.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_desc_transport_air', 'летит в воздухе') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_desc_transport_bike', 'едет на велосипеде') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_desc_transport_car', 'едет на автомобиле') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_desc_transport_unknown', 'остается на месте') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_desc_transport_walk', 'передвигается пешком') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_desc_tube', 'использует детектор [[tube]].') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_detector', 'Детектор') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_device_id', 'ID датчика') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_device_label', 'Прибор') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_extra_intro', 'Окружение') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_history_error', 'Не удалось загрузить историю.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_humidity', 'Влажность') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_last_seen', 'Последний замер') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_local_time', 'Местное время сейчас') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_no_data', 'За выбранный период данных нет.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_pressure', 'Атмосферное давление') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_temperature', 'Температура воздуха') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_temperature_f', 'Температура воздуха (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_title', 'Онлайн-датчик Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_transport', 'Передвижение') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_transport_air', 'Летательный аппарат') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_transport_bike', 'Велосипед или самокат') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_transport_car', 'Автомобиль') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_transport_unknown', 'Не обнаружено') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'live_marker_transport_walk', 'Пешком') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'locate_button_tooltip', 'Центрировать карту по моему местоположению') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'location_error', 'Произошла ошибка при получении местоположения.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'location_permission_denied', 'Доступ к местоположению отклонён.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'location_timeout', 'Время запроса местоположения истекло.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'location_unavailable', 'Местоположение недоступно.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'no_results_found', 'Результатов не найдено') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'processing_complete', 'Обработка завершена!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'processing_on_server', 'Обработка на сервере...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'qr_button_tooltip', 'QR-код ссылки на этот участок карты.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'radiation_dose', 'Мощность дозы') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'search_error', 'Ошибка поиска. Попробуйте ещё раз.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'search_input_tooltip', 'Ищите место, набирая первые буквы названия. Появится список предложений.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'select_files', 'Пожалуйста, выберите хотя бы один файл') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'short_link_tooltip', 'Нажмите, чтобы скопировать короткую ссылку') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'sources_full', 'Мы благодарим каждого, кто делится измерениями.

Людей кто анонимно загружает треки на наш сервис.
<a href="https://safecast.org" target="_blank">Safecast</a> ведёт мировую летопись радиационных данных.
<a href="https://atomfast.net" target="_blank">Atomfast</a> дозиметры AtomFast и облако Atomcloud.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> развивает линейку дозиметров-спектрометров Radiacode.

Каждый вклад расширяет общую картину; присоединяйтесь и вы.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'sources_title', 'Источники данных') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'speed', 'Скорость') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'speed_filter_tooltip_accuracy', 'Более медленные измерения остаются ближе к земле, поэтому данные пешком самые точные.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'speed_filter_tooltip_car', 'Авто: 7–200 км/ч для поездок и мобильных измерений.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'speed_filter_tooltip_intro', 'Выберите, какие измерения показывать в зависимости от скорости движения.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'speed_filter_tooltip_live', 'Сердце Safecast: данные в реальном времени с safecast.org. Используйте переключатель, чтобы показать или скрыть живые измерения.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'speed_filter_tooltip_ped', 'Пешком: менее 7 км/ч для пешеходных или стационарных измерений.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'speed_filter_tooltip_plane', 'Самолёт: 200+ км/ч для воздушных обследований.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'speed_filter_tooltip_title', 'Фильтры скорости') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'theme_toggle_tooltip', 'Переключайте светлую и тёмную темы карты.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'title', 'Чичина карта изотопов — радиоактивная карта планеты Земля') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'track_id', 'ID трека') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'upload_button', 'Загрузить [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'upload_button_tooltip', 'Добавьте свой трек замеров на карту. Поддерживаются форматы: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Можно загрузить несколько файлов, после чего откроется страница трека.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'upload_error', 'Ошибка') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'waiting_for_server', 'Ожидание сервера...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('ru', 'your_location', 'Ваше местоположение') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_archive_desc', 'Hämtar ett tgz-paket med alla publicerade .json-filer när JSON-arkivet är aktiverat.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_archive_link', 'Ladda ner veckans arkiv') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_archive_note', 'Om arkivet är avstängt returnerar servern HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_archive_title', 'Veckovis arkivpaket') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_latest_desc', 'Returnerar de senaste punkterna nära angiven latitud, longitud och sökradie i meter.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_latest_link', 'Senaste nära Tokyo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_latest_note', 'Justera lat, lon och radius_m för ditt område.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_latest_title', 'Senaste mätningarna i närheten') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_root_desc', 'Listar metadata, datasetstatistik och länkar till alla andra ändpunkter.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_root_link', 'Öppna /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_root_note', 'Börja här för att upptäcka samlingar och serverstatus.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_root_title', 'Maskinläsbart index') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_track_index_desc', 'Hämtar det N:e spåret och returnerar samma JSON som /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_track_index_link', 'Hämta index 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_track_index_note', 'Byt talet för att hämta en annan post. Svaren strömmas som radseparerad JSON.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_track_index_title', 'Ladda ner spår via index') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_tracks_desc', 'Strömmar en katalog med publicerade spår inklusive namn, antal och nedladdningslänkar.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_tracks_link', 'Visa första spåren') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_tracks_note', 'Använd parametrarna limit och offset för att bläddra i långa listor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_example_tracks_title', 'Sammanfattningar för alla spår') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_examples_heading', 'Utvalda ändpunkter') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_examples_note', 'Alla länkar svarar med JSON. När ändpunkter strömmar data kan webbläsaren ladda ner stora filer.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_intro', 'JSON-API:et speglar data som visas på kartan. Varje exempel öppnas i en ny flik så att du kan inspektera råsvaret.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_more_docs', 'Behöver du mer utförlig dokumentation?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_more_docs_link_label', 'Öppna den detaljerade guiden') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'api_title', 'Snabbstart för API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'attr_legal', 'Juridisk information') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'attr_license', 'Licens') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'attr_sources', 'Datakällor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'back_to_all_tracks', 'Tillbaka till den kombinerade spårkartan.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'count_rate', 'Räknehastighet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'date_slider_tooltip', 'Filtrera mätningar efter datum eller år.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'date_time', 'Datum och tid') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'description', 'Safecasts strålningskarta') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'download_track_cim', 'Ladda ned spår (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'duration_days', '[[count]] dagar') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'duration_hours', '[[count]] tim') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'duration_minutes', '[[count]] min') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'duration_months', '[[count]] månader') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'duration_weeks', '[[count]] veckor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'error_during_upload', 'Fel vid uppladdning!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'error_processing_files', 'Fel vid bearbetning av filer!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'files_uploaded', 'Filer uppladdade') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'geolocation_not_supported', 'Din webbläsare stöder inte geolokalisering.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'github_link_tooltip_desc', 'Ett öppet källkodsprojekt som underhålls av communityt.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'github_link_tooltip_title', 'Safecast New Map på GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'github_link_tooltip_version', 'Aktuell version: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'home_explore_global', 'Hoppa över och utforska världskartan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'home_or', 'eller') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'home_search_placeholder', 'Ange en stad, region eller land...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'home_subtitle', 'Ange en plats för att börja.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'home_title', 'Safecast strålningskarta') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'home_use_location', 'Använd min plats') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'legal_contact', 'För synpunkter kontakta:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'legal_full', 'Förord. Vi bygger en öppen karta där människor från hela världen delar dosimetermätningar för det gemensamma bästa – för vetenskap, miljö, utbildning och säkerhet. När du publicerar dina data hjälper du många. Vi ber dig ta hand om det här gemensamma arbetet.

1) Ansvar. Du ansvarar själv för riktighet och innehåll i den information du skickar in. Uppgifterna publiceras och används “som de är”. Tjänsten bekräftar eller garanterar inte att de är korrekta, fullständiga eller lämpliga för något särskilt syfte och tar inget ansvar för följderna av deras användning.

2) Öppenhet och licens. När du delar mätningar, datum, ungefärliga platser, enhetsmodeller eller andra uppgifter förstår du att de blir tillgängliga för alla och kan användas fritt under datalicensen CC0 1.0 (Public Domain). Koden finns fortsatt under MIT-licens. Du behåller upphovsrätten; ingen ersättning utgår; vidare spridning från tredje part ligger utanför vår kontroll.

3) “I befintligt skick” och utan förhandsgranskning. Bidrag publiceras utan föregående kontroll. Vi kan inte garantera instrumentkalibrering eller att uppgifterna är felfria. Informationen delas för forskningsändamål och utgör inte professionell rådgivning.

4) Integritet och moderering. För att bevara säkerhet och förtroende kan vi göra tidsangivelser och koordinater mindre precisa och ta bort eller anonymisera teknisk metadata. Vi kan dölja eller radera material som enligt vår rimliga bedömning är spam, förfalskat, olagligt eller stör tjänsten. Mätningar som skickas i god tro behandlar vi varsamt och försöker bevara.

5) Cookies. Webbplatsen använder endast en kortlivad teknisk sessionscookie; den försvinner när besöket är slut. Vi sparar inga andra spår.

Vänner, den här kartan är resultatet av gemensamma ansträngningar och öppna hjärtan. Se den som en skiss över landskapet, inte som en exakt ritning. Om vårt arbete känns meningsfullt för dig, var med – tillsammans gör vi den ännu bättre.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'legal_title', 'Juridisk information') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'legend_attention', 'Observera') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'legend_button_tooltip', 'Öppna förklaringen över strålningsnivåer.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'legend_danger', 'Fara') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'legend_full_sv', 'Denna skala visar hur säker en plats är för liv, vatten och mat.
Kom ihåg: mätningar kan vara ofullständiga; använd den bara som vägledning.

Grön (0–11 µR/h)
Nästan naturlig bakgrund.
• Brunnsvatten oftast säkert.
• Växter kan odlas utan tester.

Gul (11–30 µR/h)
Förhöjd bakgrund.
• Kontrollera vatten och jord.
• Undersök all mat innan du äter.

Röd (30–100 µR/h)
Allvarlig förorening.
• Drick inte vattnet.
• Odling eller konsumtion här är riskabelt; labbtester krävs.

Svart (>100 µR/h)
Kritisk zon.
• Vatten och mat oanvändbara.
• Endast korta besök med skydd.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'legend_safe', 'Säker') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'legend_title', 'Förklaring') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'license_full', 'Detta projekt verkar under <a href="/LICENSE" target="_blank">MIT License</a>. Hela texten finns i repots rot och på MIT:s webbplats. Du får studera, dela och förändra koden så länge dessa friheter följer med ditt arbete. Forskningsdata publiceras under licensen Creative Commons 1.0 så att mätningarna förblir i den offentliga domänen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'license_title', 'Licens') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_chart_all', 'Alla mätningar') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_chart_averaged', 'Genomsnitt över [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_chart_close', 'Stäng') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_chart_day', 'Senaste 24 timmarna') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_chart_month', 'Senaste 30 dagarna') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_chart_link', 'Öppna strålningsdiagram') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_country', 'Land') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_desc_generic', 'Denna Safecast-sensor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_desc_location', 'rapporterar från [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_desc_location_unknown', 'rapporterar från ett okänt område') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_desc_named', 'Enhet [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_desc_no_tube', 'övervakar strålningsnivåer.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_desc_radiation_sensor', 'Detta är en strålningssensor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_desc_transport_air', 'under flygning') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_desc_transport_bike', 'på cykel') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_desc_transport_car', 'i bil') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_desc_transport_unknown', 'stilla') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_desc_transport_walk', 'till fots') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_desc_tube', 'med [[tube]]-detektor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_detector', 'Detektor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_device_id', 'Enhets-ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_device_label', 'Enhet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_extra_intro', 'Miljö') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_history_error', 'Historiken kan inte laddas just nu.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_humidity', 'Luftfuktighet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_last_seen', 'Senaste mätning') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_local_time', 'Lokal tid nu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_no_data', 'Inga data registrerade under denna period.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_pressure', 'Lufttryck') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_temperature', 'Lufttemperatur') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_temperature_f', 'Lufttemperatur (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_title', 'Safecast realtidssensor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_transport', 'Rörelse') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_transport_air', 'Flygplan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_transport_bike', 'Cykel eller sparkcykel') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_transport_car', 'Bil eller skåpbil') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_transport_unknown', 'Ej upptäckt') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'live_marker_transport_walk', 'Till fots') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'locate_button_tooltip', 'Centrera kartan på min plats') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'location_error', 'Ett fel uppstod när platsen hämtades.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'location_permission_denied', 'Åtkomst till plats nekades.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'location_timeout', 'Tidsgräns för platsförfrågan överskreds.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'location_unavailable', 'Plats ej tillgänglig.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'no_results_found', 'Inga resultat hittades') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'processing_complete', 'Bearbetning klar!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'processing_on_server', 'Bearbetar på servern...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'qr_button_tooltip', 'QR-kod för länken till det här kartområdet.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'radiation_dose', 'Doshastighet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'search_error', 'Sökfel. Försök igen.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'search_input_tooltip', 'Sök efter en plats genom att skriva de första bokstäverna i namnet. En förslagslista visas.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'select_files', 'Välj minst en fil') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'short_link_tooltip', 'Klicka för att kopiera en kort delningslänk') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'sources_full', 'Vi tackar alla som delar mätningar.

Anonyma uppladdningar ritar tysta spår på kartan.
<a href="https://safecast.org" target="_blank">Safecast</a> vårdar ett globalt arkiv av avläsningar.
<a href="https://atomfast.net" target="_blank">Atomfast</a> håller Atomcloud vid liv.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> samlar insikter från Radiacode.

Varje bidrag vidgar den gemensamma bilden; vi välkomnar varmt även ditt.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'sources_title', 'Datakällor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'speed', 'Hastighet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'speed_filter_tooltip_accuracy', 'Långsammare mätningar ligger närmast marken, därför är data till fots mest precisa.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'speed_filter_tooltip_car', 'Bil: 7–200 km/h för körningar och mobila mätningar.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'speed_filter_tooltip_intro', 'Välj vilka mätningar som ska visas efter rörelsehastighet.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'speed_filter_tooltip_live', 'Safecast-hjärta: realtidsdata från safecast.org. Använd kryssrutan för att visa eller dölja live-mätningar.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'speed_filter_tooltip_ped', 'Till fots: under 7 km/h för mätningar vid gång eller stillastående.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'speed_filter_tooltip_plane', 'Flyg: 200+ km/h för mätningar från luften.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'speed_filter_tooltip_title', 'Hastighetsfilter') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'theme_toggle_tooltip', 'Växla mellan ljust och mörkt karttema.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'title', 'Safecasts isotopkarta — Radiologisk karta över planeten Jorden') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'track_id', 'Spår-ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'upload_button', 'Ladda upp [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'upload_button_tooltip', 'Lägg till din mätspårning på kartan. Format som stöds: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Du kan ladda upp flera filer och efter överföringen öppnas spårsidan.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'upload_error', 'Fel') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'waiting_for_server', 'Väntar på servern...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('sv', 'your_location', 'Din plats') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_archive_desc', 'ดาวน์โหลดแพ็กเกจ tgz ที่มีไฟล์ .json ทั้งหมดเมื่อเปิดใช้งานคลัง JSON') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_archive_link', 'ดาวน์โหลดคลังรายสัปดาห์') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_archive_note', 'หากปิดคลัง เซิร์ฟเวอร์จะส่งกลับ HTTP 404 Not Found') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_archive_title', 'แพ็กเกจคลังข้อมูลรายสัปดาห์') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_latest_desc', 'ส่งคืนจุดข้อมูลล่าสุดใกล้ละติจูด ลองจิจูด และรัศมีที่ระบุเป็นเมตร') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_latest_link', 'ล่าสุดใกล้โตเกียว') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_latest_note', 'ปรับค่า lat, lon และ radius_m ให้เหมาะกับพื้นที่ของคุณ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_latest_title', 'การวัดล่าสุดใกล้คุณ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_root_desc', 'แสดงรายการเมตาดาต้า สถิติชุดข้อมูล และลิงก์ไปยังปลายทางอื่นทั้งหมด') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_root_link', 'เปิด /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_root_note', 'เริ่มจากที่นี่เพื่อค้นหาคอลเลกชันและสถานะเซิร์ฟเวอร์') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_root_title', 'ดัชนีที่เครื่องอ่านได้') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_track_index_desc', 'ดึงเส้นทางลำดับที่ N และส่งคืน JSON แบบเดียวกับ /api/track/{id}.json') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_track_index_link', 'ดึงดัชนี 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_track_index_note', 'เปลี่ยนตัวเลขเพื่อเรียกรายการอื่น การตอบสนองถูกสตรีมเป็น JSON แบบบรรทัดต่อบรรทัด') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_track_index_title', 'ดาวน์โหลดเส้นทางตามดัชนี') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_tracks_desc', 'สตรีมแคตตาล็อกเส้นทางที่เผยแพร่พร้อมชื่อ จำนวน และลิงก์ดาวน์โหลด') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_tracks_link', 'แสดงเส้นทางชุดแรก') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_tracks_note', 'ใช้พารามิเตอร์ limit และ offset เพื่อแบ่งรายชื่อยาว ๆ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_example_tracks_title', 'สรุปเส้นทางทั้งหมด') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_examples_heading', 'ปลายทางแนะนำ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_examples_note', 'ลิงก์ทั้งหมดตอบกลับเป็น JSON เมื่อปลายทางสตรีมข้อมูล เบราว์เซอร์อาจดาวน์โหลดไฟล์ขนาดใหญ่') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_intro', 'API JSON สะท้อนข้อมูลที่แสดงบนแผนที่ ตัวอย่างแต่ละรายการจะเปิดในแท็บใหม่เพื่อให้คุณตรวจสอบผลตอบสนองดิบได้') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_more_docs', 'ต้องการคำอธิบายเชิงลึกเพิ่มเติมหรือไม่?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_more_docs_link_label', 'เปิดคู่มือฉบับละเอียด') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'api_title', 'เริ่มต้นใช้งาน API อย่างรวดเร็ว') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'attr_legal', 'ข้อมูลทางกฎหมาย') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'attr_license', 'สัญญาอนุญาต') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'attr_sources', 'แหล่งข้อมูล') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'back_to_all_tracks', 'กลับไปยังแผนที่รวมเส้นทาง') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'count_rate', 'อัตราการนับ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'date_slider_tooltip', 'กรองค่าการวัดตามวันที่หรือช่วงปี.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'date_time', 'วันที่และเวลา') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'description', 'แผนที่รังสีของ Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'download_track_cim', 'ดาวน์โหลดเส้นทาง (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'duration_days', '[[count]] วัน') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'duration_hours', '[[count]] ชม.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'duration_minutes', '[[count]] นาที') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'duration_months', '[[count]] เดือน') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'duration_weeks', '[[count]] สัปดาห์') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'error_during_upload', 'เกิดข้อผิดพลาดระหว่างการอัปโหลด!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'error_processing_files', 'เกิดข้อผิดพลาดขณะประมวลผลไฟล์!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'files_uploaded', 'อัปโหลดไฟล์แล้ว') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'geolocation_not_supported', 'เบราว์เซอร์ของคุณไม่รองรับการระบุตำแหน่งทางภูมิศาสตร์') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'github_link_tooltip_desc', 'โครงการโอเพนซอร์สที่ชุมชนดูแลร่วมกัน.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'github_link_tooltip_title', 'Safecast New Map บน GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'github_link_tooltip_version', 'เวอร์ชันปัจจุบัน: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'home_explore_global', 'ข้ามและสำรวจแผนที่โลก') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'home_or', 'หรือ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'home_search_placeholder', 'ป้อนเมือง ภูมิภาค หรือประเทศ...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'home_subtitle', 'ป้อนสถานที่เพื่อเริ่มต้น') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'home_title', 'แผนที่รังสี Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'home_use_location', 'ใช้ตำแหน่งของฉัน') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'legal_contact', 'หากต้องการส่งความคิดเห็น โปรดติดต่อ:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'legal_full', 'คำนำ เรากำลังสร้างแผนที่แบบเปิดที่ผู้คนทั่วโลกแบ่งปันค่าที่อ่านได้จากเครื่องโดซิเมตรเพื่อประโยชน์ส่วนรวม — เพื่อวิทยาศาสตร์ สิ่งแวดล้อม การศึกษา และความปลอดภัย เมื่อคุณเผยแพร่ข้อมูล คุณได้ช่วยเหลือผู้คนมากมาย โปรดดูแลผลงานร่วมชิ้นนี้ด้วยความเอาใจใส่

1) ความรับผิดชอบ ความถูกต้องและเนื้อหาของข้อมูลที่คุณส่งเป็นความรับผิดชอบของคุณเอง ข้อมูลถูกเผยแพร่และใช้งาน “ตามสภาพ” บริการนี้ไม่ได้รับรองหรือรับประกันความถูกต้อง ความครบถ้วน หรือความเหมาะสมสำหรับวัตถุประสงค์ใด ๆ และไม่รับผิดชอบต่อผลที่อาจเกิดขึ้นจากการใช้งาน

2) ความเปิดกว้างและสัญญาอนุญาต เมื่อคุณแบ่งปันค่าที่วัดได้ วันที่ ตำแหน่งโดยประมาณ รุ่นอุปกรณ์ หรือรายละเอียดอื่น ๆ คุณเข้าใจว่าข้อมูลเหล่านี้จะเปิดให้ทุกคนเข้าถึงและนำไปใช้ได้อย่างเสรีภายใต้สัญญาอนุญาตข้อมูล CC0 1.0 (สาธารณสมบัติ) โค้ดยังคงเผยแพร่ภายใต้สัญญาอนุญาต MIT คุณยังคงเป็นเจ้าของผลงาน ไม่มีค่าตอบแทน และการนำไปเผยแพร่ต่อโดยบุคคลที่สามอยู่นอกเหนือการควบคุมของเรา

3) “ตามสภาพ” และไม่มีการตรวจสอบล่วงหน้า เนื้อหาถูกเผยแพร่โดยไม่มีการตรวจสอบก่อน เราไม่สามารถรับประกันการปรับเทียบเครื่องมือหรือการปราศจากข้อผิดพลาดได้ ข้อมูลถูกแบ่งปันเพื่อการวิจัยและไม่ใช่คำแนะนำทางวิชาชีพ

4) ความเป็นส่วนตัวและการดูแลความเรียบร้อย เพื่อรักษาความปลอดภัยและความไว้วางใจ เราอาจปรับเวลาและพิกัดให้กว้างขึ้น รวมทั้งลบหรือทำให้ข้อมูลเมตาทางเทคนิคเป็นนิรนาม เราอาจซ่อนหรือ ลบเนื้อหาที่ตามการพิจารณาอย่างสมเหตุสมผลของเราแล้วเป็นสแปม ปลอมแปลง ผิดกฎหมาย หรือรบกวนการให้บริการ เราดูแลข้อมูลการวัดที่มาด้วยความตั้งใจดีอย่างระมัดระวังและพยายามรักษาไว้

5) คุกกี้ เว็บไซต์นี้ใช้คุกกี้เซสชันทางเทคนิคอายุสั้นเพียงตัวเดียว ซึ่งจะหายไปเมื่อคุณสิ้นสุดการเยี่ยมชม เราไม่เก็บร่องรอยอื่นใดเพิ่มเติม

เพื่อน ๆ ที่รัก แผนที่นี้เกิดจากความพยายามร่วมกันและหัวใจที่เปิดกว้าง จงมองมันเป็นภาพร่างของภูมิประเทศ ไม่ใช่แผนที่ที่ละเอียดทุกเส้น หากงานของเราส่งพลังให้คุณ โปรดมาร่วมกัน — เราจะทำให้มันดียิ่งขึ้นได้ด้วยกัน') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'legal_title', 'ข้อมูลทางกฎหมาย') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'legend_attention', 'ระวัง') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'legend_button_tooltip', 'เปิดคำอธิบายระดับรังสี.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'legend_danger', 'อันตราย') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'legend_full_th', 'สเกลนี้บอกว่าพื้นที่ปลอดภัยต่อชีวิต น้ำ และอาหารแค่ไหน
จำไว้: การวัดอาจไม่ครบถ้วน ใช้เป็นเพียงแนวทาง

เขียว (0–11 µR/h)
ฉากหลังเกือบธรรมชาติ
• น้ำบ่อส่วนใหญ่ปลอดภัย
• ปลูกพืชได้โดยไม่ต้องตรวจ

เหลือง (11–30 µR/h)
ฉากหลังสูงขึ้น
• ตรวจน้ำและดิน
• ตรวจอาหารก่อนกิน

แดง (30–100 µR/h)
ปนเปื้อนรุนแรง
• ห้ามดื่มน้ำ
• ปลูกหรือกินที่นี่เสี่ยง ต้องตรวจแลบ

ดำ (>100 µR/h)
เขตวิกฤต
• น้ำและอาหารใช้ไม่ได้
• อยู่ได้เพียงสั้นๆ พร้อมอุปกรณ์ป้องกัน
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'legend_safe', 'ปลอดภัย') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'legend_title', 'คำอธิบาย') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'license_full', 'โครงการนี้ดำเนินงานภายใต้สัญญาอนุญาต <a href="/LICENSE" target="_blank">MIT License</a> ข้อความฉบับเต็มอยู่ที่รากของรีโพซิทอรีและบนเว็บไซต์ MIT คุณสามารถศึกษา แชร์ และปรับแก้โค้ดได้ ตราบใดที่เสรีภาพเหล่านี้ติดตามผลงานของคุณไปด้วย ข้อมูลการวิจัยเผยแพร่ภายใต้สัญญาอนุญาต Creative Commons 1.0 เพื่อให้ค่าการวัดยังคงเป็นสาธารณะ.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'license_title', 'สัญญาอนุญาต') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_chart_all', 'การอ่านทั้งหมด') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_chart_averaged', 'เฉลี่ยใน [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_chart_close', 'ปิด') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_chart_day', '24 ชั่วโมงล่าสุด') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_chart_month', '30 วันล่าสุด') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_chart_link', 'เปิดกราฟรังสี') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_country', 'ประเทศ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_desc_generic', 'เซ็นเซอร์ Safecast นี้') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_desc_location', 'รายงานจาก [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_desc_location_unknown', 'รายงานจากพื้นที่ไม่ทราบ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_desc_named', 'อุปกรณ์ [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_desc_no_tube', 'เฝ้าระวังระดับรังสี') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_desc_radiation_sensor', 'นี่คือเซ็นเซอร์รังสี') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_desc_transport_air', 'ขณะบิน') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_desc_transport_bike', 'ด้วยจักรยาน') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_desc_transport_car', 'ด้วยรถยนต์') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_desc_transport_unknown', 'อยู่กับที่') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_desc_transport_walk', 'เดินเท้า') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_desc_tube', 'ด้วยตัวตรวจจับ [[tube]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_detector', 'ตัวตรวจจับ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_device_id', 'รหัสอุปกรณ์') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_device_label', 'อุปกรณ์') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_extra_intro', 'สิ่งแวดล้อม') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_history_error', 'ไม่สามารถโหลดประวัติได้ในขณะนี้') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_humidity', 'ความชื้น') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_last_seen', 'การอ่านล่าสุด') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_local_time', 'เวลาท้องถิ่นปัจจุบัน') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_no_data', 'ไม่มีข้อมูลในช่วงเวลานี้') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_pressure', 'ความดันอากาศ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_temperature', 'อุณหภูมิอากาศ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_temperature_f', 'อุณหภูมิอากาศ (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_title', 'เซ็นเซอร์ Safecast แบบเรียลไทม์') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_transport', 'การเคลื่อนที่') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_transport_air', 'เครื่องบิน') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_transport_bike', 'จักรยานหรือสกูตเตอร์') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_transport_car', 'รถยนต์หรือรถตู้') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_transport_unknown', 'ไม่พบ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'live_marker_transport_walk', 'เดินเท้า') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'locate_button_tooltip', 'จัดกึ่งกลางแผนที่ไปยังตำแหน่งของฉัน') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'location_error', 'เกิดข้อผิดพลาดระหว่างรับตำแหน่ง') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'location_permission_denied', 'ปฏิเสธการเข้าถึงตำแหน่ง') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'location_timeout', 'หมดเวลาคำขอตำแหน่ง') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'location_unavailable', 'ไม่สามารถใช้ตำแหน่งได้') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'no_results_found', 'ไม่พบผลลัพธ์') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'processing_complete', 'ประมวลผลเสร็จสิ้น!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'processing_on_server', 'กำลังประมวลผลบนเซิร์ฟเวอร์...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'qr_button_tooltip', 'คิวอาร์โค้ดของลิงก์สำหรับพื้นที่แผนที่นี้') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'radiation_dose', 'อัตราปริมาณรังสี') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'search_error', 'ข้อผิดพลาดในการค้นหา กรุณาลองอีกครั้ง') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'search_input_tooltip', 'ค้นหาสถานที่โดยพิมพ์ตัวอักษรแรกของชื่อ รายการแนะนำจะปรากฏขึ้น') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'select_files', 'โปรดเลือกอย่างน้อยหนึ่งไฟล์') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'short_link_tooltip', 'คลิกเพื่อคัดลอกลิงก์สั้นสำหรับแชร์') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'sources_full', 'เราขอขอบคุณทุกคนที่แบ่งปันการวัดค่า

การอัปโหลดแบบไม่ระบุชื่อวาดเส้นทางเงียบ ๆ บนแผนที่
<a href="https://safecast.org" target="_blank">Safecast</a> ดูแลคลังข้อมูลการอ่านระดับโลก
<a href="https://atomfast.net" target="_blank">Atomfast</a> รักษาให้ Atomcloud สว่างอยู่เสมอ
<a href="https://radiaverse.com" target="_blank">Radioverse</a> รวบรวมองค์ความรู้จาก Radiacode

ทุกการมีส่วนร่วมช่วยขยายภาพรวมร่วมกัน เราขอเชิญคุณมาร่วมเติมข้อมูลของคุณด้วยนะ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'sources_title', 'แหล่งข้อมูล') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'speed', 'ความเร็ว') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'speed_filter_tooltip_accuracy', 'การวัดที่ช้ากว่าจะอยู่ใกล้พื้นที่สุด ดังนั้นข้อมูลการเดินเท้าจึงแม่นยำที่สุด.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'speed_filter_tooltip_car', 'รถยนต์: 7–200 กม./ชม. สำหรับการขับขี่และการวัดแบบเคลื่อนที่.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'speed_filter_tooltip_intro', 'เลือกการวัดที่จะให้แสดงตามความเร็วในการเคลื่อนที่.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'speed_filter_tooltip_live', 'หัวใจ Safecast: ข้อมูลเรียลไทม์จาก safecast.org ใช้ตัวเลือกนี้เพื่อแสดงหรือซ่อนค่าการวัดแบบสด') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'speed_filter_tooltip_ped', 'เดินเท้า: ต่ำกว่า 7 กม./ชม. สำหรับการวัดขณะเดินหรือหยุดนิ่ง.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'speed_filter_tooltip_plane', 'เครื่องบิน: 200+ กม./ชม. สำหรับการสำรวจทางอากาศ.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'speed_filter_tooltip_title', 'ตัวกรองความเร็ว') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'theme_toggle_tooltip', 'สลับระหว่างธีมแผนที่แบบสว่างและมืด.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'title', 'แผนที่ไอโซโทปของ Safecast — แผนที่รังสีของโลก') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'track_id', 'รหัสเส้นทาง') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'upload_button', 'อัปโหลด [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'upload_button_tooltip', 'เพิ่มเส้นทางการวัดของคุณลงในแผนที่ รูปแบบที่รองรับ: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log คุณสามารถอัปโหลดหลายไฟล์ได้ หลังจากอัปโหลดแล้วจะเปิดหน้าของเส้นทาง') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'upload_error', 'ข้อผิดพลาด') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'waiting_for_server', 'กำลังรอเซิร์ฟเวอร์...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('th', 'your_location', 'ตำแหน่งของคุณ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_archive_desc', 'JSON arşivi etkinse yayımlanmış tüm .json dosyalarıyla bir tgz paketi indirir.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_archive_link', 'Haftalık arşivi indir') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_archive_note', 'Arşiv devre dışıysa sunucu HTTP 404 Not Found döndürür.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_archive_title', 'Haftalık arşiv paketi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_latest_desc', 'Belirtilen enlem, boylam ve arama yarıçapı (metre) yakınındaki en yeni noktaları döndürür.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_latest_link', 'Tokyo yakınındaki en yeniler') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_latest_note', 'Kendi bölgeniz için lat, lon ve radius_m değerlerini ayarlayın.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_latest_title', 'Yakındaki en yeni ölçümler') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_root_desc', 'Meta verileri, veri kümesi istatistiklerini ve diğer tüm uç noktalara bağlantıları listeler.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_root_link', '/api’yi aç') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_root_note', 'Koleksiyonları ve sunucu durumunu keşfetmek için buradan başlayın.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_root_title', 'Makinece okunabilir dizin') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_track_index_desc', 'N’inci rotayı çözümler ve /api/track/{id}.json ile aynı JSON’u döndürür.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_track_index_link', '1. indeksi getir') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_track_index_note', 'Farklı bir kayıt almak için sayıyı değiştirin. Yanıtlar satır satır JSON olarak iletilir.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_track_index_title', 'Dizine göre rota indir') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_tracks_desc', 'Yayınlanan rotaların adlarını, sayımlarını ve indirme bağlantılarını içeren bir kataloğu yayınlar.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_tracks_link', 'İlk rotaları listele') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_tracks_note', 'Uzun listelerde gezinmek için limit ve offset parametrelerini kullanın.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_example_tracks_title', 'Tüm rota özetleri') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_examples_heading', 'Öne çıkan uç noktalar') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_examples_note', 'Tüm bağlantılar JSON döndürür. Uç noktalar veri akıttığında tarayıcı büyük dosyalar indirebilir.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_intro', 'JSON API, haritada gösterilen verileri yansıtır. Her örnek yeni bir sekmede açılır, böylece ham yanıtı inceleyebilirsiniz.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_more_docs', 'Daha ayrıntılı belgelere mi ihtiyacınız var?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_more_docs_link_label', 'Ayrıntılı rehberi aç') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'api_title', 'API hızlı başlangıç') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'attr_legal', 'Hukuki bilgiler') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'attr_license', 'Lisans') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'attr_sources', 'Veri kaynakları') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'back_to_all_tracks', 'Tüm izlerin birleşik haritasına geri dön.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'count_rate', 'Sayım hızı') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'date_slider_tooltip', 'Ölçümleri tarihe veya yıl aralığına göre filtreleyin.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'date_time', 'Tarih ve saat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'description', 'Safecast’nın radyasyon haritası') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'download_track_cim', 'Parkuru indir (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'duration_days', '[[count]] gün') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'duration_hours', '[[count]] sa') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'duration_minutes', '[[count]] dk') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'duration_months', '[[count]] ay') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'duration_weeks', '[[count]] hafta') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'error_during_upload', 'Yükleme sırasında hata oluştu!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'error_processing_files', 'Dosyalar işlenirken hata oluştu!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'files_uploaded', 'Yüklenen dosyalar') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'geolocation_not_supported', 'Tarayıcınız jeolokasyonu desteklemiyor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'github_link_tooltip_desc', 'Topluluk tarafından yönetilen açık kaynaklı proje.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'github_link_tooltip_title', 'GitHub''da Safecast New Map') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'github_link_tooltip_version', 'Güncel sürüm: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'home_explore_global', 'Atla ve dünya haritasını keşfet') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'home_or', 'veya') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'home_search_placeholder', 'Şehir, bölge veya ülke girin...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'home_subtitle', 'Başlamak için bir konum girin.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'home_title', 'Safecast radyolojik harita') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'home_use_location', 'Konumumu kullan') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'legal_contact', 'Geri bildirim için yazın:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'legal_full', 'Önsöz. Bilimin, çevrenin, eğitimin ve güvenliğin ortak yararı için dünyanın dört bir yanından insanların dozimetre ölçümlerini paylaştığı açık bir harita geliştiriyoruz. Verilerini yayımladığında pek çok kişiye destek oluyorsun. Lütfen bu ortak emeğe özen göster.

1) Sorumluluk. Gönderdiğin bilgilerin doğruluğu ve içeriği senin sorumluluğundadır. Veriler “olduğu gibi” yayımlanır ve kullanılır. Hizmet, verilerin doğruluğunu, eksiksizliğini ya da herhangi bir amaç için uygunluğunu onaylamaz veya garanti etmez; kullanımından doğabilecek sonuçlardan sorumluluk kabul etmez.

2) Açıklık ve lisans. Ölçümler, tarihler, yaklaşık konumlar, cihaz modelleri veya diğer ayrıntıları paylaştığında, bunların herkesin erişimine açılacağını ve CC0 1.0 (Kamusal Alan) veri lisansı kapsamında serbestçe kullanılabileceğini kabul edersin. Kod MIT lisansı altında erişilebilir kalır. Eser sahipliği sende kalır; ödeme yapılmaz; üçüncü tarafların yeniden dağıtımı kontrolümüz dışındadır.

3) “Olduğu gibi” ve doğrulama olmadan. Paylaşımlar ön inceleme yapılmadan yayımlanır. Cihazların kalibrasyonunu veya hatasızlığını garanti edemeyiz. Bilgiler araştırma amaçlı paylaşılır ve profesyonel tavsiye niteliği taşımaz.

4) Gizlilik ve moderasyon. Güvenlik ve güven duygusunu korumak için zaman damgalarını ve koordinatları yuvarlayabilir, teknik meta verileri silebilir veya anonimleştirebiliriz. Makul değerlendirmemize göre spam, sahte, yasa dışı ya da hizmeti aksatan içerikleri gizleyebilir veya silebiliriz. İyi niyetle gönderilen ölçümleri özenle saklar ve korumaya çalışırız.

5) Çerezler. Site yalnızca kısa ömürlü bir teknik oturum çerezi kullanır; ziyaretin bittiğinde çerez de silinir. Başka izler tutmayız.

Sevgili dostlar, bu harita ortak çabanın ve açık yüreklerin ürünüdür. Onu kusursuz bir plan yerine arazinin bir eskizi olarak düşünün. Çalışmamız size anlamlı geliyorsa aramıza katılın — birlikte daha da iyisini yapabiliriz.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'legal_title', 'Yasal bilgiler') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'legend_attention', 'Dikkat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'legend_button_tooltip', 'Radyasyon seviyeleri açıklamasını açın.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'legend_danger', 'Tehlike') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'legend_full_tr', 'Bu ölçek, bir yerin yaşam, su ve gıda için ne kadar güvenli olduğunu gösterir.
Unutma: ölçümler eksik olabilir; yalnızca rehber olarak kullan.

Yeşil (0–11 µR/h)
Doğal arka plan.
• Kuyu suyu genelde güvenlidir.
• Test olmadan bitki yetiştirilebilir.

Sarı (11–30 µR/h)
Yükselmiş arka plan.
• Su ve toprağı kontrol et.
• Her gıdayı yemeden önce test et.

Kırmızı (30–100 µR/h)
Ciddi kirlilik.
• Suyu içme.
• Burada yetiştirmek veya yemek riskli; lab testleri şart.

Siyah (>100 µR/h)
Kritik bölge.
• Su ve yiyecek kullanılamaz.
• Uzun süre kalmak yasak; sadece kısa süre korumayla kal.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'legend_safe', 'Güvenli') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'legend_title', 'Lejant') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'license_full', 'Bu proje <a href="/LICENSE" target="_blank">MIT License</a> lisansı altında büyüyor. Tam metin depo kök dizininde ve MIT sitesinde bulunur. Kodunu inceleyebilir, paylaşabilir ve değiştirebilirsin; yeter ki bu özgürlükler çalışmalarına da eşlik etsin. Araştırma verileri Creative Commons 1.0 lisansı ile yayımlanır; böylece ölçümler kamu malı olarak kalır.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'license_title', 'Lisans') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_chart_all', 'Tüm okumalar') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_chart_averaged', '[[window]] ortalaması') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_chart_close', 'Kapat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_chart_day', 'Son 24 saat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_chart_month', 'Son 30 gün') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_chart_link', 'Radyasyon grafiklerini aç') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_country', 'Ülke') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_desc_generic', 'Bu Safecast sensörü') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_desc_location', '[[place]] konumundan raporluyor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_desc_location_unknown', 'bilinmeyen bir bölgeden raporluyor') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_desc_named', 'Cihaz [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_desc_no_tube', 'radyasyon seviyelerini izliyor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_desc_radiation_sensor', 'Bu bir radyasyon sensörüdür.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_desc_transport_air', 'uçuş sırasında') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_desc_transport_bike', 'bisikletle') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_desc_transport_car', 'arabayla') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_desc_transport_unknown', 'sabit') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_desc_transport_walk', 'yürüyerek') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_desc_tube', '[[tube]] dedektörü ile.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_detector', 'Dedektör') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_device_id', 'Cihaz kimliği') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_device_label', 'Cihaz') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_extra_intro', 'Çevre') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_history_error', 'Geçmiş şu anda yüklenemiyor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_humidity', 'Nem') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_last_seen', 'Son okuma') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_local_time', 'Şu anki yerel saat') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_no_data', 'Bu dönemde kayıtlı veri yok.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_pressure', 'Hava basıncı') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_temperature', 'Hava sıcaklığı') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_temperature_f', 'Hava sıcaklığı (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_title', 'Safecast gerçek zamanlı sensör') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_transport', 'Hareket') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_transport_air', 'Uçak') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_transport_bike', 'Bisiklet veya scooter') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_transport_car', 'Araba veya minibüs') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_transport_unknown', 'Algılanmadı') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'live_marker_transport_walk', 'Yürüyerek') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'locate_button_tooltip', 'Haritayı konumuma ortala') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'location_error', 'Konum alınırken bir hata oluştu.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'location_permission_denied', 'Konuma erişim reddedildi.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'location_timeout', 'Konum isteği zaman aşımına uğradı.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'location_unavailable', 'Konum kullanılamıyor.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'no_results_found', 'Sonuç bulunamadı') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'processing_complete', 'İşleme tamamlandı!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'processing_on_server', 'Sunucuda işleniyor...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'qr_button_tooltip', 'Bu harita bölümü bağlantısının QR kodu.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'radiation_dose', 'Doz hızı') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'search_error', 'Arama hatası. Lütfen tekrar deneyin.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'search_input_tooltip', 'Adının ilk harflerini yazarak bir yer arayın. Öneri listesi görünecektir.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'select_files', 'Lütfen en az bir dosya seçin') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'short_link_tooltip', 'Kısa paylaşım bağlantısını kopyalamak için tıklayın') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'sources_full', 'Ölçümlerini paylaşan herkese teşekkür ederiz.

Anonim yüklemeler haritada sessiz izler çizer.
<a href="https://safecast.org" target="_blank">Safecast</a> küresel bir ölçüm arşivini yaşatır.
<a href="https://atomfast.net" target="_blank">Atomfast</a> Atomcloud''un ışığını açık tutar.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> Radiacode deneyimlerini derler.

Her katkı ortak manzarayı genişletir; kendi verini de içtenlikle bekliyoruz.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'sources_title', 'Veri kaynakları') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'speed', 'Hız') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'speed_filter_tooltip_accuracy', 'Daha yavaş ölçümler zemine en yakın kalır; bu yüzden yaya verileri en doğrusudur.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'speed_filter_tooltip_car', 'Araba: sürüşler ve mobil ölçümler için 7–200 km/s.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'speed_filter_tooltip_intro', 'Hangi ölçümlerin hareket hızına göre gösterileceğini seçin.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'speed_filter_tooltip_live', 'Safecast kalbi: safecast.org’dan gerçek zamanlı veriler. Bu seçeneği kullanarak canlı ölçümleri gösterip gizleyebilirsiniz.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'speed_filter_tooltip_ped', 'Yaya: yürürken veya dururken yapılan ölçümler için 7 km/s altı.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'speed_filter_tooltip_plane', 'Uçak: hava ölçümleri için 200+ km/s.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'speed_filter_tooltip_title', 'Hız filtreleri') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'theme_toggle_tooltip', 'Haritanın açık ve koyu teması arasında geçiş yapın.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'title', 'Safecast İzotop Haritası — Dünya gezegeninin radyolojik haritası') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'track_id', 'İz Kimliği') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'upload_button', 'Yükle [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'upload_button_tooltip', 'Ölçüm rotanızı haritaya ekleyin. Desteklenen formatlar: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Birden fazla dosya yükleyebilir, yükleme sonrası rota sayfası açılır.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'upload_error', 'Hata') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'waiting_for_server', 'Sunucu bekleniyor...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('tr', 'your_location', 'Konumunuz') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_archive_desc', 'Завантажує tgz-пакет з усіма опублікованими файлами .json, коли JSON-архів увімкнено.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_archive_link', 'Завантажити тижневий архів') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_archive_note', 'Якщо архів вимкнено, сервер поверне HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_archive_title', 'Тижневий архів') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_latest_desc', 'Повертає найновіші точки поблизу заданих широти, довготи й радіуса пошуку в метрах.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_latest_link', 'Останні біля Токіо') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_latest_note', 'Налаштуйте lat, lon і radius_m під власний район.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_latest_title', 'Останні вимірювання поруч') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_root_desc', 'Перелічує метадані, статистику набору та посилання на всі інші ендпоінти.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_root_link', 'Відкрити /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_root_note', 'Почніть тут, щоб дізнатися про колекції та стан сервера.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_root_title', 'Машиночитний індекс') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_track_index_desc', 'Повертає N-ий маршрут і той самий JSON, що й /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_track_index_link', 'Отримати індекс 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_track_index_note', 'Змініть число, щоб отримати інший запис. Відповіді передаються построчним JSON.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_track_index_title', 'Завантажити маршрут за індексом') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_tracks_desc', 'Стрімить каталог опублікованих маршрутів з назвами, кількістю та посиланнями на завантаження.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_tracks_link', 'Список перших маршрутів') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_tracks_note', 'Використовуйте параметри limit і offset, щоб розбивати довгі списки.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_example_tracks_title', 'Усі зведення маршрутів') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_examples_heading', 'Рекомендовані ендпоінти') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_examples_note', 'Усі посилання відповідають у форматі JSON. Коли кінцеві точки транслюють дані, браузер може завантажувати великі файли.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_intro', 'JSON‑API відтворює дані, які ви бачите на карті. Кожен приклад відкривається в новій вкладці, щоб переглянути сирий відгук.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_more_docs', 'Потрібно більше деталей?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_more_docs_link_label', 'Відкрити докладний посібник') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'api_title', 'Швидкий старт API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'attr_legal', 'Правова інформація') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'attr_license', 'Ліцензія') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'attr_sources', 'Джерела даних') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'back_to_all_tracks', 'Повернутися до об’єднаної карти треків.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'count_rate', 'Швидкість лічби') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'date_slider_tooltip', 'Фільтруйте вимірювання за датою або діапазоном років.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'date_time', 'Дата і час') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'description', 'Карта радіації Чічі') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'download_track_cim', 'Завантажити трек (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'duration_days', '[[count]] днів') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'duration_hours', '[[count]] год') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'duration_minutes', '[[count]] хв') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'duration_months', '[[count]] місяців') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'duration_weeks', '[[count]] тижнів') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'error_during_upload', 'Помилка під час завантаження!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'error_processing_files', 'Помилка обробки файлів!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'files_uploaded', 'Файли завантажено') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'geolocation_not_supported', 'Ваш браузер не підтримує геолокацію.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'github_link_tooltip_desc', 'Проєкт з відкритим кодом, який підтримує спільнота.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'github_link_tooltip_title', 'Safecast New Map на GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'github_link_tooltip_version', 'Поточна версія: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'home_explore_global', 'Пропустити й дослідити світову карту') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'home_or', 'або') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'home_search_placeholder', 'Введіть місто, регіон або країну...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'home_subtitle', 'Введіть місце для початку.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'home_title', 'Радіологічна карта Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'home_use_location', 'Використати моє місцезнаходження') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'legal_contact', 'Для відгуків пишіть на:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'legal_full', 'Преамбула. Ми створюємо відкриту мапу, де люди з усього світу діляться показниками дозиметрів для спільного блага — науки, довкілля, освіти й безпеки. Публікуючи свої дані, ви допомагаєте багатьом. Просимо дбайливо ставитися до цієї спільної справи.

1) Відповідальність. Точність і зміст інформації, яку ви надсилаєте, залишаються вашою відповідальністю. Дані публікуються та використовуються «як є». Сервіс не підтверджує і не гарантує їхню правильність, повноту чи придатність для будь-якої мети та не несе відповідальності за можливі наслідки їх використання.

2) Відкритість і ліцензія. Поширюючи вимірювання, дати, приблизні координати, модель приладу чи інші відомості, ви розумієте, що вони стають доступними для всіх і можуть вільно використовуватися за ліцензією CC0 1.0 (Public Domain) для даних. Код проєкту залишається під ліцензією MIT. Авторство зберігається за вами; винагорода не передбачена; подальше розповсюдження третіми сторонами не контролюється сервісом.

3) «Як є» і без перевірки. Матеріали публікуються без попереднього перегляду. Ми не можемо гарантувати повірку приладів або відсутність помилок. Інформація надається для дослідницьких цілей і не є професійною рекомендацією.

4) Приватність і модерація. Заради безпеки та довіри ми можемо узагальнювати час і координати, видаляти або анонімізувати технічні метадані. Ми маємо право приховувати чи видаляти матеріали, які, на нашу розумну думку, є спамом, підробкою, порушують закон або заважають роботі сервісу. До добросовісних вимірювань ставимося уважно й намагаємося їх зберегти.

5) Файли cookie. Сайт використовує лише короткочасний технічний cookie сесії; він зникає після завершення візиту. Жодних інших слідів ми не зберігаємо.

Друзі, ця мапа — результат спільної праці та відкритих сердець. Сприймайте її як ескіз місцевості, а не точний кресленик. Якщо наша робота вам відгукується, долучайтеся — разом зробимо її кращою.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'legal_title', 'Правова інформація') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'legend_attention', 'Увага') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'legend_button_tooltip', 'Відкрити легенду рівнів радіації.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'legend_danger', 'Небезпека') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'legend_full_uk', 'Ця шкала показує, наскільки місце безпечне для життя, води й їжі.
Пам’ятайте: вимірювання можуть бути неповними; користуйтеся лише як орієнтиром.

Зелений (0–11 µR/h)
Майже природний фон.
• Кринична вода здебільшого безпечна.
• Рослини можна вирощувати без перевірок.

Жовтий (11–30 µR/h)
Підвищений фон.
• Перевіряйте воду й ґрунт.
• Тестуйте будь-яку їжу перед споживанням.

Червоний (30–100 µR/h)
Серйозне забруднення.
• Воду пити не можна.
• Вирощування чи вживання продуктів тут небезпечно; потрібні лабораторні дослідження.

Чорний (>100 µR/h)
Критична зона.
• Вода й їжа непридатні.
• Лише коротке перебування з захистом.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'legend_safe', 'Безпечно') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'legend_title', 'Легенда') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'license_full', 'Цей проєкт розвивається під ліцензією <a href="/LICENSE" target="_blank">MIT License</a>. Повний текст лежить у корені репозиторію та на сайті MIT. Ви можете вивчати, поширювати й змінювати код, якщо ці свободи залишаються разом із вашою роботою. Дані досліджень публікуються за ліцензією Creative Commons 1.0, тож вимірювання залишаються у суспільному надбанні.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'license_title', 'Ліцензія') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_chart_all', 'Усі показники') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_chart_averaged', 'Середнє за [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_chart_close', 'Закрити') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_chart_day', 'Останні 24 години') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_chart_month', 'Останні 30 днів') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_chart_link', 'Відкрити графіки радіації') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_country', 'Країна') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_desc_generic', 'Цей сенсор Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_desc_location', 'звітує з [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_desc_location_unknown', 'звітує з невідомої області') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_desc_named', 'Пристрій [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_desc_no_tube', 'відстежує рівні радіації.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_desc_radiation_sensor', 'Це сенсор радіації.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_desc_transport_air', 'під час польоту') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_desc_transport_bike', 'на велосипеді') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_desc_transport_car', 'на автомобілі') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_desc_transport_unknown', 'нерухомий') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_desc_transport_walk', 'пішки') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_desc_tube', 'з детектором [[tube]].') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_detector', 'Детектор') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_device_id', 'ID пристрою') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_device_label', 'Пристрій') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_extra_intro', 'Середовище') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_history_error', 'Наразі не вдається завантажити історію.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_humidity', 'Вологість') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_last_seen', 'Останній показник') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_local_time', 'Місцевий час зараз') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_no_data', 'За цей період дані не зафіксовано.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_pressure', 'Атмосферний тиск') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_temperature', 'Температура повітря') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_temperature_f', 'Температура повітря (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_title', 'Сенсор Safecast у реальному часі') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_transport', 'Рух') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_transport_air', 'Літак') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_transport_bike', 'Велосипед або самокат') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_transport_car', 'Автомобіль або фургон') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_transport_unknown', 'Не виявлено') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'live_marker_transport_walk', 'Пішки') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'locate_button_tooltip', 'Центрувати карту за моїм розташуванням') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'location_error', 'Під час отримання розташування сталася помилка.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'location_permission_denied', 'Доступ до розташування відхилено.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'location_timeout', 'Час запиту розташування вичерпано.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'location_unavailable', 'Розташування недоступне.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'no_results_found', 'Результатів не знайдено') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'processing_complete', 'Обробку завершено!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'processing_on_server', 'Обробка на сервері...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'qr_button_tooltip', 'QR-код посилання на цю ділянку карти.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'radiation_dose', 'Потужність дози') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'search_error', 'Помилка пошуку. Спробуйте ще раз.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'search_input_tooltip', 'Шукайте місце, набираючи перші літери назви. З''явиться список пропозицій.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'select_files', 'Будь ласка, виберіть щонайменше один файл') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'short_link_tooltip', 'Натисніть, щоб скопіювати коротке посилання для поширення') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'sources_full', 'Дякуємо всім, хто ділиться вимірюваннями.

Анонімні завантаження малюють на мапі тихі стежки.
<a href="https://safecast.org" target="_blank">Safecast</a> підтримує світовий архів показників.
<a href="https://atomfast.net" target="_blank">Atomfast</a> тримає Atomcloud у строю.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> збирає напрацювання Radiacode.

Кожен внесок розширює спільну картину; щиро запрошуємо додати й ваші дані.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'sources_title', 'Джерела даних') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'speed', 'Швидкість') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'speed_filter_tooltip_accuracy', 'Повільніші вимірювання залишаються найближче до землі, тому пішохідні дані найточніші.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'speed_filter_tooltip_car', 'Авто: 7–200 км/год для поїздок і мобільних вимірювань.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'speed_filter_tooltip_intro', 'Виберіть, які вимірювання показувати за швидкістю руху.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'speed_filter_tooltip_live', 'Серце Safecast: дані в реальному часі із safecast.org. Перемикайте опцію, щоб показувати або ховати живі вимірювання.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'speed_filter_tooltip_ped', 'Пішки: менше 7 км/год для пішохідних чи стаціонарних вимірювань.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'speed_filter_tooltip_plane', 'Літак: 200+ км/год для аерозйомок.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'speed_filter_tooltip_title', 'Фільтри швидкості') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'theme_toggle_tooltip', 'Перемикайте світлу й темну тему карти.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'title', 'Карта ізотопів Чічі — Радіологічна карта планети Земля') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'track_id', 'ID треку') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'upload_button', 'Завантажити [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'upload_button_tooltip', 'Додайте свій трек вимірювань на карту. Підтримувані формати: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Можна завантажити кілька файлів, після чого відкриється сторінка треку.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'upload_error', 'Помилка') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'waiting_for_server', 'Очікування сервера...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('uk', 'your_location', 'Ваше розташування') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_archive_desc', 'Tải gói tgz chứa mọi tệp .json đã công bố khi bật lưu trữ JSON.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_archive_link', 'Tải lưu trữ hàng tuần') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_archive_note', 'Nếu lưu trữ bị tắt, máy chủ sẽ trả về HTTP 404 Not Found.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_archive_title', 'Gói lưu trữ hàng tuần') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_latest_desc', 'Trả về các điểm mới nhất gần vĩ độ, kinh độ và bán kính tìm kiếm đã cho (tính bằng mét).') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_latest_link', 'Mới nhất gần Tokyo') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_latest_note', 'Điều chỉnh lat, lon và radius_m cho khu vực của bạn.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_latest_title', 'Đo đạc mới nhất gần bạn') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_root_desc', 'Liệt kê siêu dữ liệu, thống kê bộ dữ liệu và liên kết đến mọi điểm cuối khác.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_root_link', 'Mở /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_root_note', 'Bắt đầu tại đây để khám phá các bộ sưu tập và trạng thái máy chủ.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_root_title', 'Chỉ mục đọc được bằng máy') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_track_index_desc', 'Truy xuất hành trình thứ N và trả về cùng tài liệu JSON như /api/track/{id}.json.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_track_index_link', 'Lấy chỉ số 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_track_index_note', 'Thay số khác để lấy mục khác. Phản hồi được truyền ở dạng JSON từng dòng.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_track_index_title', 'Tải hành trình theo chỉ số') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_tracks_desc', 'Phát danh mục các hành trình đã công bố cùng tên, số lượng và liên kết tải xuống.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_tracks_link', 'Liệt kê các hành trình đầu tiên') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_tracks_note', 'Dùng tham số limit và offset để phân trang danh sách dài.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_example_tracks_title', 'Tổng quan mọi hành trình') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_examples_heading', 'Điểm cuối nổi bật') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_examples_note', 'Tất cả liên kết đều trả về JSON. Khi đầu cuối phát trực tuyến dữ liệu, trình duyệt có thể tải xuống tệp lớn.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_intro', 'API JSON phản chiếu dữ liệu hiển thị trên bản đồ. Mỗi ví dụ mở trong một thẻ mới để bạn xem phản hồi thô.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_more_docs', 'Cần tài liệu chi tiết hơn?') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_more_docs_link_label', 'Mở hướng dẫn chi tiết') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'api_title', 'Khởi động nhanh API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'attr_legal', 'Thông tin pháp lý') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'attr_license', 'Giấy phép') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'attr_sources', 'Nguồn dữ liệu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'back_to_all_tracks', 'Quay lại bản đồ tổng hợp các hành trình.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'count_rate', 'Tốc độ đếm') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'date_slider_tooltip', 'Lọc số liệu theo ngày hoặc khoảng năm.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'date_time', 'Ngày và giờ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'description', 'Bản đồ bức xạ của Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'download_track_cim', 'Tải bản ghi (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'duration_days', '[[count]] ngày') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'duration_hours', '[[count]] giờ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'duration_minutes', '[[count]] phút') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'duration_months', '[[count]] tháng') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'duration_weeks', '[[count]] tuần') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'error_during_upload', 'Lỗi trong quá trình tải lên!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'error_processing_files', 'Lỗi khi xử lý tệp!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'files_uploaded', 'Tệp đã được tải lên') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'geolocation_not_supported', 'Trình duyệt của bạn không hỗ trợ định vị địa lý.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'github_link_tooltip_desc', 'Dự án mã nguồn mở do cộng đồng duy trì.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'github_link_tooltip_title', 'Safecast New Map trên GitHub') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'github_link_tooltip_version', 'Phiên bản hiện tại: {version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'home_explore_global', 'Bỏ qua và khám phá bản đồ thế giới') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'home_or', 'hoặc') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'home_search_placeholder', 'Nhập thành phố, vùng hoặc quốc gia...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'home_subtitle', 'Nhập một địa điểm để bắt đầu.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'home_title', 'Bản đồ phóng xạ Safecast') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'home_use_location', 'Sử dụng vị trí của tôi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'legal_contact', 'Để góp ý, hãy liên hệ:') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'legal_full', 'Lời mở đầu. Chúng tôi đang xây dựng một bản đồ mở nơi mọi người trên khắp thế giới chia sẻ số liệu từ máy đo liều phóng xạ vì lợi ích chung — cho khoa học, môi trường, giáo dục và an toàn. Khi bạn công bố dữ liệu, bạn giúp được rất nhiều người. Mong bạn trân trọng nỗ lực chung này.

1) Trách nhiệm. Bạn chịu trách nhiệm về độ chính xác và nội dung của thông tin do chính bạn gửi lên. Dữ liệu được công bố và sử dụng “nguyên trạng”. Dịch vụ không xác nhận hay bảo đảm tính đúng đắn, đầy đủ hoặc sự phù hợp của dữ liệu cho bất kỳ mục đích nào, và không chịu trách nhiệm cho những hệ quả phát sinh từ việc sử dụng dữ liệu đó.

2) Tính mở và giấy phép. Khi chia sẻ các phép đo, thời gian, vị trí ước chừng, mẫu thiết bị hoặc thông tin khác, bạn hiểu rằng tất cả sẽ được công khai và có thể được sử dụng tự do theo giấy phép dữ liệu CC0 1.0 (Phạm vi công cộng). Mã nguồn tiếp tục được phát hành theo giấy phép MIT. Quyền tác giả vẫn thuộc về bạn; chúng tôi không chi trả thù lao; việc phân phối lại bởi bên thứ ba nằm ngoài khả năng kiểm soát của chúng tôi.

3) “Nguyên trạng” và chưa được kiểm chứng. Nội dung được đăng tải mà không qua kiểm duyệt trước. Chúng tôi không thể bảo đảm việc hiệu chuẩn của thiết bị hoặc việc không có sai sót. Thông tin được chia sẻ nhằm mục đích nghiên cứu và không phải là khuyến nghị chuyên môn.

4) Quyền riêng tư và điều phối. Để giữ an toàn và xây dựng niềm tin, chúng tôi có thể làm mờ thời gian và tọa độ, đồng thời gỡ bỏ hoặc ẩn danh siêu dữ liệu kỹ thuật. Chúng tôi có thể ẩn hoặc xóa những nội dung mà theo đánh giá hợp lý của mình là spam, giả mạo, trái pháp luật hoặc gây rối dịch vụ. Các số liệu được gửi với thiện chí sẽ được chúng tôi nâng niu và cố gắng lưu giữ.

5) Cookie. Trang web chỉ dùng một cookie phiên làm việc kỹ thuật tồn tại trong thời gian ngắn; cookie sẽ biến mất khi bạn kết thúc lượt truy cập. Chúng tôi không lưu giữ thêm dấu vết nào khác.

Bạn hữu thân mến, bản đồ này là kết quả của nỗ lực chung và những tấm lòng rộng mở. Hãy xem đây như một phác thảo địa hình chứ không phải bản vẽ chính xác tuyệt đối. Nếu cảm thấy đồng điệu với công việc của chúng tôi, mời bạn cùng tham gia — chúng ta sẽ cùng nhau làm cho nó tốt hơn nữa.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'legal_title', 'Thông tin pháp lý') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'legend_attention', 'Chú ý') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'legend_button_tooltip', 'Mở chú giải mức bức xạ.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'legend_danger', 'Nguy hiểm') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'legend_full_vi', 'Thang này cho biết mức an toàn của nơi ở đối với đời sống, nước và thực phẩm.
Nhớ rằng: đo đạc có thể chưa đầy đủ; chỉ dùng để tham khảo.

Xanh (0–11 µR/h)
Nền gần tự nhiên.
• Nước giếng thường an toàn.
• Có thể trồng trọt không cần thử nghiệm.

Vàng (11–30 µR/h)
Nền cao.
• Kiểm tra nước và đất.
• Thử mọi thực phẩm trước khi ăn.

Đỏ (30–100 µR/h)
Ô nhiễm nghiêm trọng.
• Không uống nước.
• Trồng hoặc ăn ở đây nguy hiểm; cần xét nghiệm.

Đen (>100 µR/h)
Vùng cực nguy.
• Nước và thức ăn không dùng được.
• Chỉ ở lại ngắn với bảo hộ.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'legend_safe', 'An toàn') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'legend_title', 'Chú thích') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'license_full', 'Dự án này phát triển dưới <a href="/LICENSE" target="_blank">MIT License</a>. Toàn văn nằm ở thư mục gốc của kho và trên trang MIT. Bạn có thể nghiên cứu, chia sẻ và chỉnh sửa mã miễn là những quyền tự do này theo sát sản phẩm của bạn. Dữ liệu nghiên cứu được phát hành theo giấy phép Creative Commons 1.0 để các phép đo luôn thuộc phạm vi công cộng.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'license_title', 'Giấy phép') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_chart_all', 'Tất cả số đọc') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_chart_averaged', 'Trung bình trong [[window]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_chart_close', 'Đóng') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_chart_day', '24 giờ qua') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_chart_month', '30 ngày qua') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_chart_link', 'Mở biểu đồ bức xạ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_country', 'Quốc gia') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_desc_generic', 'Cảm biến Safecast này') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_desc_location', 'báo cáo từ [[place]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_desc_location_unknown', 'báo cáo từ khu vực không xác định') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_desc_named', 'Thiết bị [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_desc_no_tube', 'giám sát mức bức xạ.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_desc_radiation_sensor', 'Đây là cảm biến bức xạ.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_desc_transport_air', 'khi bay') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_desc_transport_bike', 'bằng xe đạp') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_desc_transport_car', 'bằng ô tô') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_desc_transport_unknown', 'đứng yên') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_desc_transport_walk', 'đi bộ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_desc_tube', 'với đầu dò [[tube]].') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_detector', 'Đầu dò') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_device_id', 'ID thiết bị') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_device_label', 'Thiết bị') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_extra_intro', 'Môi trường') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_history_error', 'Không thể tải lịch sử lúc này.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_humidity', 'Độ ẩm') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_last_seen', 'Số đọc mới nhất') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_local_time', 'Giờ địa phương hiện tại') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_no_data', 'Không có dữ liệu trong giai đoạn này.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_pressure', 'Áp suất không khí') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_temperature', 'Nhiệt độ không khí') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_temperature_f', 'Nhiệt độ không khí (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_title', 'Cảm biến Safecast thời gian thực') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_transport', 'Di chuyển') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_transport_air', 'Máy bay') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_transport_bike', 'Xe đạp hoặc xe tay ga') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_transport_car', 'Ô tô hoặc xe tải nhỏ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_transport_unknown', 'Không phát hiện') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'live_marker_transport_walk', 'Đi bộ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'locate_button_tooltip', 'Căn giữa bản đồ theo vị trí của tôi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'location_error', 'Đã xảy ra lỗi khi lấy vị trí.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'location_permission_denied', 'Bị từ chối quyền truy cập vị trí.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'location_timeout', 'Hết thời gian yêu cầu vị trí.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'location_unavailable', 'Không thể xác định vị trí.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'no_results_found', 'Không tìm thấy kết quả') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'processing_complete', 'Xử lý xong!') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'processing_on_server', 'Đang xử lý trên máy chủ...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'qr_button_tooltip', 'Mã QR của liên kết đến khu vực bản đồ này.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'radiation_dose', 'Suất liều') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'search_error', 'Lỗi tìm kiếm. Vui lòng thử lại.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'search_input_tooltip', 'Tìm kiếm địa điểm bằng cách nhập các chữ cái đầu tiên của tên. Danh sách gợi ý sẽ xuất hiện.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'select_files', 'Vui lòng chọn ít nhất một tệp') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'short_link_tooltip', 'Nhấn để sao chép liên kết chia sẻ ngắn') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'sources_full', 'Chúng tôi cảm ơn tất cả những ai chia sẻ phép đo.

Các lượt tải ẩn danh vẽ nên những đường mòn lặng lẽ trên bản đồ.
<a href="https://safecast.org" target="_blank">Safecast</a> nuôi dưỡng kho lưu trữ số liệu toàn cầu.
<a href="https://atomfast.net" target="_blank">Atomfast</a> giữ cho Atomcloud luôn sáng.
<a href="https://radiaverse.com" target="_blank">Radioverse</a> gom góp hiểu biết từ Radiacode.

Mỗi đóng góp đều mở rộng bức tranh chung; rất mong bạn cũng góp phần của mình.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'sources_title', 'Nguồn dữ liệu') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'speed', 'Tốc độ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'speed_filter_tooltip_accuracy', 'Các phép đo chậm sát mặt đất hơn, vì vậy dữ liệu đi bộ là chính xác nhất.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'speed_filter_tooltip_car', 'Ô tô: 7–200 km/h cho các chuyến đi và phép đo di động.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'speed_filter_tooltip_intro', 'Chọn phép đo được hiển thị theo tốc độ di chuyển.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'speed_filter_tooltip_live', 'Trái tim Safecast: dữ liệu thời gian thực từ safecast.org. Dùng lựa chọn này để hiện hoặc ẩn các phép đo trực tiếp.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'speed_filter_tooltip_ped', 'Đi bộ: dưới 7 km/h cho phép đo khi đi bộ hoặc đứng yên.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'speed_filter_tooltip_plane', 'Máy bay: 200+ km/h cho khảo sát trên không.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'speed_filter_tooltip_title', 'Bộ lọc tốc độ') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'theme_toggle_tooltip', 'Chuyển đổi giữa chủ đề bản đồ sáng và tối.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'title', 'Bản đồ đồng vị của Safecast — Bản đồ phóng xạ của Trái Đất') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'track_id', 'ID hành trình') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'upload_button', 'Tải lên [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'upload_button_tooltip', 'Thêm hành trình đo đạc của bạn vào bản đồ. Định dạng hỗ trợ: .kml, .kmz, .gpx, .csv, .rctrk, .json, .log. Bạn có thể tải lên nhiều tệp; sau khi tải lên sẽ mở trang hành trình.') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'upload_error', 'Lỗi') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'waiting_for_server', 'Đang chờ máy chủ...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('vi', 'your_location', 'Vị trí của bạn') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_archive_desc', '在启用 JSON 归档时，下载包含所有已发布 .json 文件的 tgz 包。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_archive_link', '下载每周归档') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_archive_note', '如果归档被禁用，服务器将返回 HTTP 404 Not Found。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_archive_title', '每周归档包') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_latest_desc', '返回靠近指定纬度、经度和搜索半径（米）的最新点。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_latest_link', '东京附近的最新数据') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_latest_note', '调整 lat、lon 和 radius_m 以聚焦你的区域。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_latest_title', '附近的最新测量') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_root_desc', '列出元数据、数据集统计以及指向其他所有端点的链接。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_root_link', '打开 /api') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_root_note', '从这里开始了解集合和服务器状态。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_root_title', '机器可读索引') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_track_index_desc', '解析第 N 条轨迹并返回与 /api/track/{id}.json 相同的 JSON。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_track_index_link', '获取索引 1') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_track_index_note', '更改数字以获取其他条目。响应以按行分隔的 JSON 形式流式传输。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_track_index_title', '按索引下载轨迹') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_tracks_desc', '以流方式返回已发布轨迹的目录，包括名称、数量和下载链接。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_tracks_link', '列出首批轨迹') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_tracks_note', '使用 limit 和 offset 参数浏览长列表。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_example_tracks_title', '所有轨迹摘要') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_examples_heading', '推荐的端点') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_examples_note', '所有链接都返回 JSON。当端点以流形式提供数据时，浏览器可能会下载大文件。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_intro', 'JSON API 会反映地图上显示的数据。每个示例都会在新标签页打开，方便你查看原始响应。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_method_get', 'GET') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_more_docs', '需要更详细的说明？') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_more_docs_link_label', '打开详细指南') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'api_title', 'API 快速入门') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'attr_api', 'API') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'attr_legal', '法律信息') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'attr_license', '许可') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'attr_sources', '数据来源') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'back_to_all_tracks', '返回到合并轨迹的地图。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'count_rate', '计数率') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'date_slider_tooltip', '按日期或年份范围筛选测量数据。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'date_time', '日期和时间') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'description', 'Safecast 的辐射地图') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'download_track_cim', '下载轨迹 (.json)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'duration_days', '[[count]]天') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'duration_hours', '[[count]]小时') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'duration_minutes', '[[count]]分钟') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'duration_months', '[[count]]个月') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'duration_weeks', '[[count]]周') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'error_during_upload', '上传过程中出错！') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'error_processing_files', '处理文件时出错！') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'files_uploaded', '文件已上传') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'geolocation_not_supported', '您的浏览器不支持地理定位。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'github_link_tooltip_desc', '由社区维护的开源项目。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'github_link_tooltip_title', 'GitHub 上的 Safecast New Map') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'github_link_tooltip_version', '当前版本：{version}') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'home_explore_global', '跳过并浏览全球地图') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'home_or', '或') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'home_search_placeholder', '输入城市、地区或国家...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'home_subtitle', '输入一个地点开始。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'home_title', 'Safecast辐射地图') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'home_use_location', '使用我的位置') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'legal_contact', '如需反馈，欢迎联系：') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'legal_full', '前言。我们正在构建一张开放地图，全球各地的人们在此分享辐射剂量计读数，服务于科学、生态、教育和安全等公共利益。发布您的数据，将帮助许多人。请善待这份共同的成果。

1) 责任。您提交的任何信息的准确性和内容由您自行负责。数据按"原样"发布和使用。本服务不确认或保证其正确性、完整性或对任何目的的适用性，也不对使用这些数据可能产生的任何后果承担责任。

2) 开放与许可。通过分享剂量测量数据、日期、大致位置、设备型号或其他事实信息，您理解这些内容将对所有人可用，且数据可在 CC0 1.0（公共领域）许可下自由使用，代码则在 MIT 许可下提供。著作权归您所有；不提供报酬；第三方的进一步传播超出我们的控制范围。

3) "原样"且未经验证。发布的内容未经事先审核。我们无法保证仪器校准的准确性或数据中没有错误。这些信息用于研究目的分享，不构成专业建议。

4) 隐私与审核。为保持地图的安全性和可信度，时间戳和坐标可能会被概括化处理，技术元数据可能被移除或匿名化。我们可能隐藏或删除在我们合理判断下属于垃圾信息、伪造、违法或有破坏性的内容。我们尊重真实的测量数据，并努力保存它们。

5) Cookies。本站仅使用短期技术性会话 Cookie；它在您的访问结束后即消失。我们不保留任何其他痕迹。

朋友们，这张地图是共同努力和开放心灵的结晶。请将它视为地形的草图，而非精确的蓝图。如果我们的工作引起了您的共鸣，请加入我们——一起让它变得更好。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'legal_title', '法律须知') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'legend_attention', '注意') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'legend_button_tooltip', '打开辐射等级图例。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'legend_danger', '危险') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'legend_full_en', 'This scale shows how likely a spot is safe for folk, water, and food.
Mind: the readings might not be complete, and some rays could be higher or hiding. Treat these numbers as guidance only.

Green (0–11 µR/h)
Background near natural.
• Water from wells is generally safe.
• You can grow plants without checks.

Yellow (11–30 µR/h)
Raised background; take care.
• Check water and soil.
• Test any veg, mushrooms, or other produce before you eat.

Red (30–100 µR/h)
Serious contamination.
• Don''t drink the water.
• Growing or eating produce from here is risky; lab tests are a must.

Black (>100 µR/h)
Critical zone.
• Water and food cannot be used.
• Staying long-term is out; only short visits with protection.
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'legend_full_zh', '此刻度展示一个地点对生命、水和食物的安全程度。
请记住：读数可能不完整；仅作参考。

绿色 (0–11 µR/h)
背景接近自然。
• 井水通常安全。
• 可放心种植。

黄色 (11–30 µR/h)
背景升高。
• 检查水和土。
• 食物食用前需检测。

红色 (30–100 µR/h)
严重污染。
• 不可饮用水。
• 种植或食用此地产品危险；需实验室检测。

黑色 (>100 µR/h)
危急区域。
• 水和食物不可用。
• 仅短暂停留且需防护。
') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'legend_safe', '安全') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'legend_title', '图例') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'license_full', '本项目基于 <a href="/LICENSE" target="_blank">MIT 许可证</a> 发布。完整文本位于代码仓库根目录及 MIT 官网。您可以学习、分享和修改代码，前提是这些自由权利随您的作品一同传递。研究数据集以 Creative Commons 1.0 (CC0) 发布，确保测量数据保留在公共领域。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'license_title', '许可') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_chart_all', '所有读数') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_chart_averaged', '[[window]]平均值') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_chart_close', '关闭') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_chart_day', '过去24小时') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_chart_month', '过去30天') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_chart_link', '打开辐射图表') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_country', '国家') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_desc_generic', '此Safecast传感器') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_desc_location', '从[[place]]报告') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_desc_location_unknown', '从未知区域报告') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_desc_named', '设备 [[name]]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_desc_no_tube', '正在监测辐射水平。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_desc_radiation_sensor', '这是一个辐射传感器。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_desc_transport_air', '飞行中') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_desc_transport_bike', '骑自行车中') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_desc_transport_car', '驾车中') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_desc_transport_unknown', '静止') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_desc_transport_walk', '步行中') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_desc_tube', '使用[[tube]]探测器。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_detector', '探测器') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_device_id', '设备ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_device_label', '设备') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_extra_intro', '环境') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_history_error', '目前无法加载历史记录。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_humidity', '湿度') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_last_seen', '最新读数') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_local_time', '当地时间') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_no_data', '此时段无记录数据。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_pressure', '气压') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_temperature', '气温') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_temperature_f', '气温 (°F)') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_title', 'Safecast实时传感器') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_transport', '移动方式') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_transport_air', '飞机') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_transport_bike', '自行车或滑板车') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_transport_car', '汽车或面包车') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_transport_unknown', '未检测到') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'live_marker_transport_walk', '步行') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'locate_button_tooltip', '将地图居中到我的位置') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'location_error', '获取位置时发生错误。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'location_permission_denied', '位置访问被拒绝。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'location_timeout', '位置请求超时。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'location_unavailable', '位置不可用。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'no_results_found', '未找到结果') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'processing_complete', '处理完成！') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'processing_on_server', '服务器处理中...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'qr_button_tooltip', '此地图区域链接的二维码。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'radiation_dose', '剂量率') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'search_error', '搜索错误，请重试。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'search_input_tooltip', '输入地名的前几个字符来搜索。将显示建议列表。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'select_files', '请至少选择一个文件') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'short_link_tooltip', '点击复制简短的分享链接') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'sources_full', '感谢所有分享测量数据的人。

匿名上传在地图上留下静默的轨迹。
<a href="https://safecast.org" target="_blank">Safecast</a> 维护着一个全球测量数据档案。
<a href="https://atomfast.net" target="_blank">Atomfast</a> 保持着 Atomcloud 的运行。
<a href="https://radiaverse.com" target="_blank">Radioverse</a> 收集来自 Radiacode 的见解。

每一份贡献都在拓宽共同的图景；我们热忱地邀请您加入。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'sources_title', '数据来源') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'speed', '速度') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'speed_filter_tooltip_accuracy', '速度越慢的测量越贴近地面，因此步行数据最精确。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'speed_filter_tooltip_car', '汽车：7–200 公里/小时，用于驾车和移动测量。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'speed_filter_tooltip_intro', '根据移动速度选择要显示的测量数据。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'speed_filter_tooltip_live', 'Safecast 心形：来自 safecast.org 的实时数据。使用该选项即可显示或隐藏实时测量。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'speed_filter_tooltip_ped', '步行：低于 7 公里/小时，用于步行或静止测量。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'speed_filter_tooltip_plane', '飞机：200+ 公里/小时，用于空中测量。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'speed_filter_tooltip_title', '速度筛选') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'theme_toggle_tooltip', '在地图的亮色和暗色主题之间切换。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'title', 'Safecast 同位素地图——地球辐射地图') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'track_id', '轨迹 ID') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'upload_button', '上传 [+]') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'upload_button_tooltip', '将您的测量轨迹添加到地图。支持的格式：.kml、.kmz、.gpx、.csv、.rctrk、.json、.log。您可以上传多个文件，上传完成后将打开轨迹页面。') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'upload_error', '错误') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'waiting_for_server', '正在等待服务器...') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();
INSERT INTO translations (language_code, key, value) VALUES ('zh', 'your_location', '您的位置') ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW();

COMMIT;