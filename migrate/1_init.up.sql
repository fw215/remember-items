
--
-- テーブルの構造 `categories`
--

CREATE TABLE `categories` (
  `category_id` bigint(20) NOT NULL COMMENT 'カテゴリID',
  `user_id` bigint(20) NOT NULL COMMENT 'ユーザーID',
  `category_name` text NOT NULL COMMENT 'カテゴリ名',
  `created` datetime NOT NULL COMMENT '登録日時',
  `modified` datetime NOT NULL COMMENT '更新日時'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='カテゴリ';

-- --------------------------------------------------------

--
-- テーブルの構造 `items`
--

CREATE TABLE `items` (
  `item_id` bigint(20) NOT NULL COMMENT 'アイテムID',
  `user_id` bigint(20) NOT NULL COMMENT 'ユーザーID',
  `category_id` bigint(20) NOT NULL COMMENT 'カテゴリID',
  `item_name` text NOT NULL COMMENT 'アイテム名',
  `item_image` text COMMENT 'アイテム画像',
  `created` datetime NOT NULL COMMENT '登録日時',
  `modified` datetime NOT NULL COMMENT '更新日時'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='アイテム';

-- --------------------------------------------------------

--
-- テーブルの構造 `users`
--

CREATE TABLE `users` (
  `user_id` bigint(20) NOT NULL COMMENT 'ユーザーID',
  `google_id` varchar(255) DEFAULT NULL COMMENT 'GoogleID',
  `google_email` varchar(255) DEFAULT NULL COMMENT 'Googleメールアドレス',
  `google_access_token` varchar(255) DEFAULT NULL COMMENT 'Googleアクセストークン',
  `google_expiry` datetime DEFAULT NULL COMMENT 'Google有効期限',
  `google_refresh_token` varchar(255) NOT NULL COMMENT 'googleリフレッシュトークン',
  `created` datetime NOT NULL COMMENT '登録日時',
  `modified` datetime NOT NULL COMMENT '更新日時'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='ユーザー';

--
-- Indexes for dumped tables
--

--
-- Indexes for table `categories`
--
ALTER TABLE `categories`
  ADD PRIMARY KEY (`category_id`);

--
-- Indexes for table `items`
--
ALTER TABLE `items`
  ADD PRIMARY KEY (`item_id`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`user_id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `categories`
--
ALTER TABLE `categories`
  MODIFY `category_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'カテゴリID';
--
-- AUTO_INCREMENT for table `items`
--
ALTER TABLE `items`
  MODIFY `item_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'アイテムID';
--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `user_id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'ユーザーID';