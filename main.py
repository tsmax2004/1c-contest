from os import listdir
from os.path import isfile, join
import diff_match_patch as dmp_module


path1 = input()
path2 = input()
match_border = float(input())

files1 = [join(path1, f) for f in listdir(path1) if isfile(join(path1, f))]
unique_files1 = set(files1)
files2 = [join(path2, f) for f in listdir(path2) if isfile(join(path2, f))]
unique_files2 = set(files2)

identical_files = []
similar_files = []

for file_path1 in files1:
    with open(file_path1, 'r') as f1:
        content1 = f1.read()
        for file_path2 in files2:
            with open(file_path2, 'r') as f2:
                content2 = f2.read()

                if len(content1) < len(content2):
                    content1, content2 = content2, content1

                dmp = dmp_module.diff_match_patch()
                diff = dmp.diff_main(content1, content2)

                common = 0
                for batch in diff:
                    if batch[0] == 0:
                        common += len(batch[1])

                match = common / len(content1)

                if match == 1:
                    identical_files.append((file_path1, file_path2))
                    unique_files1.discard(file_path1)
                    unique_files2.discard(file_path2)
                elif match >= match_border:
                    similar_files.append((file_path1, file_path2, match))
                    unique_files1.discard(file_path1)
                    unique_files2.discard(file_path2)

print("Идентичные файлы:")
for f1, f2 in identical_files:
    print(f1, f2)
print()

print("Похожие файлы:")
for f1, f2, match in similar_files:
    print(f1, f2, match)
print()

print(f'Уникальные файлы в {path1}:')
print('\n'.join(unique_files1))
print()

print(f'Уникальные файлы в {path2}:')
print('\n'.join(unique_files2))