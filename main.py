import sys, getopt
from os import listdir
from os.path import isfile, isdir, join
import diff_match_patch as dmp_module


class DirComparator:
    def __init__(self, dir_path1: str, dir_path2: str, match_border: int, levenshtein: bool = False):
        self.dir_path1 = dir_path1
        self.dir_path2 = dir_path2
        self.match_border = match_border
        self.levenshtein = levenshtein

        self.identical_files: list[tuple[str, str]] = []
        self.similar_files: list[tuple[str, str, float]] = []
        self.unique_files1: set[str] = set()
        self.unique_files2: set[str] = set()

    def compute_match_percent(self, diff, max_len: int) -> float:
        common = 0
        for batch in diff:
            if batch[0] == 0:
                common += len(batch[1])
        return common / max_len * 100

    def compare_contents(self, content1: str, content2: str) -> float:
        dmp = dmp_module.diff_match_patch()
        diff = dmp.diff_main(content1, content2)
        max_len = max(len(content1), len(content2))

        if self.levenshtein:
            return (max_len - dmp.diff_levenshtein(diff)) / max_len * 100
        return self.compute_match_percent(diff, max_len)

    def print_results(self):
        if len(self.identical_files) != 0:
            print("Идентичные файлы:")
            for f1, f2 in self.identical_files:
                print(f1, f2)
            print()

        if len(self.similar_files) != 0:
            print("Похожие файлы:")
            for f1, f2, m in self.similar_files:
                print(f'{f1} {f2} - {m:.2f}% сходства')
            print()

        if len(self.unique_files1) != 0:
            print("Уникальные файлы в первой директории:")
            print('\n'.join(self.unique_files1))
            print()

        if len(self.unique_files2) != 0:
            print("Уникальные файлы во второй директории:")
            print('\n'.join(self.unique_files2))
            print()

    def compare_dirs(self):
        files1 = [join(self.dir_path1, f) for f in listdir(self.dir_path1) if isfile(join(self.dir_path1, f))]
        if len(files1) == 0:
            print("Первая директория пустая")
            return

        files2 = [join(self.dir_path2, f) for f in listdir(self.dir_path2) if isfile(join(self.dir_path2, f))]
        if len(files2) == 0:
            print("Вторая директория пустая")
            return

        self.unique_files1 = set(files1)
        self.unique_files2 = set(files2)

        self.identical_files.clear()
        self.similar_files.clear()

        for file_path1 in files1:
            with open(file_path1, 'r') as f1:
                content1 = f1.read()
                for file_path2 in files2:
                    with open(file_path2, 'r') as f2:
                        content2 = f2.read()

                        if len(content1) == 0 and len(content2) == 0:
                            match = 1
                        elif len(content1) == 0 or len(content2) == 0:
                            continue
                        elif len(content1) / len(content2) * 100 < self.match_border \
                                or len(content2) / len(content1) * 100 < self.match_border:
                            continue
                        else:
                            match = self.compare_contents(content1, content2)

                        if match == 100:
                            self.identical_files.append((file_path1, file_path2))
                            self.unique_files1.discard(file_path1)
                            self.unique_files2.discard(file_path2)
                        elif match >= self.match_border:
                            self.similar_files.append((file_path1, file_path2, match))
                            self.unique_files1.discard(file_path1)
                            self.unique_files2.discard(file_path2)

        self.print_results()


def main(argv):
    levenshtein = False
    try:
        opts, args = getopt.getopt(argv, "hl")
    except getopt.GetoptError as e:
        print(f'Неизвестный параметр: {e.opt}')
        print('Использование: main.py [-l] dir1 dir2 match')
        print('Попробуйте -h для подробной информации')
        return

    for opt, arg in opts:
        if opt == '-h':
            print('Использование: main.py [-l] dir1 dir2 match')
            print('-l         - использовать расстояние Левенштейна')
            print('dir1, dir2 - пути до директорий')
            print('match      - требуемый процент сходства (целое число от 0 до 100)')
            return
        elif opt == '-l':
            levenshtein = True

    if len(args) == 0:
        dir_path1 = input("Введите абсолютный путь первой директории: ")
        dir_path2 = input("Введите абсолютный путь второй директории: ")

        try:
            match_border = int(input("Введите требуемый процент сходства: "))
        except TypeError:
            print("Неверный ввод: процент сходства должен быть целым числом от 0 до 100")
            return
        print()

    elif len(args) == 3:
        dir_path1, dir_path2, match_border = args
        try:
            match_border = int(match_border)
        except TypeError:
            print("Неверный ввод: процент сходства должен быть целым числом от 0 до 100")
            return
    else:
        print(f'Ожидалось 3 аргумента (путь1, путь2, процент сходства), получено {len(args)}')
        return

    if not isdir(dir_path1):
        print("Некорректный путь к первой директории")
        return

    if not isdir(dir_path2):
        print("Некорректный путь ко второй директории")
        return

    if dir_path1 == dir_path2:
        print("Директории совпадают")
        return

    if not (0 <= match_border <= 100):
        print("Неверный ввод: процент сходства должен быть целым числом от 0 до 100")
        return

    cmp = DirComparator(dir_path1, dir_path2, match_border, levenshtein=levenshtein)
    cmp.compare_dirs()


if __name__ == '__main__':
    main(sys.argv[1:])
