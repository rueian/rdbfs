#! /bin/sh

test_cmp()
{
	echo "test content $1 $2"
	cmp $1 $2
	if [ $? -eq 0 ]; then
		echo "PASS: same content"
	else
		echo "FAIL: content diff"
	fi
}

symlink_test()
{
	mkdir testlink
	sync
	sleep 2
	echo "test content" > testlink/the_file
	sync
	sleep .5
	ln -s testlink/the_file sym
	sync
	sleep .5

	if [ -e sym ]; then
		echo "PASS: sym link create success"
	else
		echo "FAIL: sym link create failed"
	fi

	test_cmp sym testlink/the_file
	
	rm testlink/the_file
	cat sym
	if [ $? -eq 0 ]; then
		echo "FAIL: sym still linked"
	else
		echo "PASS: link to removed file"
	fi

	rm sym
	rm -r testlink
}

hardlink_test()
{
	echo 'content' > file1
	ln file1 file2
	ln file2 file3
	test_cmp file1 file2
	test_cmp file3 file2
	test_cmp file3 file1

	echo "new data" >> file1
	test_cmp file1 file2
	test_cmp file3 file2
	test_cmp file3 file1
	echo "new data222" >> file2
	test_cmp file1 file2
	test_cmp file3 file2
	test_cmp file3 file1

	rm file1
	test_cmp file3 file2
	echo "voaiwejfoi" >> file3
	test_cmp file3 file2

	rm file2 file3
}

symlink_test
hardlink_test

